package githubsource

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gritcli/grit/daemon/internal/logs"
	"github.com/gritcli/grit/daemon/internal/statuspage"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const (
	clientID     = "976b90cbef967ca64b7e"
	clientSecret = "83abf75db10bf0d38f8cc6db263f7835cc8943fd"
)

// SignIn signs in to the source.
func (s *source) SignIn(
	ctx context.Context,
	log logs.Log,
) error {
	if s.config.Token != "" {
		return errors.New("already authenticated using a personal access token (PAT)")
	}

	if isEnterpriseServer(s.config.Domain) {
		return errors.New("sign-in to GitHub Enterprise Server is not supported, use a personal access token (PAT) instead")
	}

	lis, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("cannot start oauth callback listener: %w", err)
	}
	defer lis.Close()

	log.WriteVerbose("listening for oauth callback on %s", lis.Addr())

	cfg := s.oauthConfig(lis.Addr())
	state := uuid.NewString()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	svr := &http.Server{}
	svr.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			c := r.URL.Query().Get("code")
			s := r.URL.Query().Get("state")

			if c == "" || s != state {
				statuspage.RenderDefault(w, r, http.StatusBadRequest)
				return
			}

			_, err := cfg.Exchange(r.Context(), c)
			if err != nil {
				statuspage.RenderDefault(w, r, http.StatusBadRequest)
				return
			}

			renderAuthSuccess(w, r)
			cancel()
		},
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := svr.Serve(lis)
		if err == http.ErrServerClosed {
			return nil
		}

		return fmt.Errorf("oauth callback server failed: %w", err)
	})

	g.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		return svr.Shutdown(ctx)
	})

	g.Go(func() error {
		url := cfg.AuthCodeURL(state)

		if err := browser.OpenURL(url); err != nil {
			return fmt.Errorf("unable to open URL in browser: %w", err)
		}

		return nil
	})

	return g.Wait()
}

// SignOut signs out of the source.
func (s *source) SignOut(
	ctx context.Context,
	log logs.Log,
) error {
	return errors.New("<not implemented>")
}

func (s *source) oauthConfig(addr net.Addr) oauth2.Config {
	u := &url.URL{
		Scheme: "http",
		Host:   addr.String(),
	}

	return oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: u.String(),
		Scopes:      []string{"repo"},
	}
}

func renderAuthSuccess(w http.ResponseWriter, r *http.Request) {
	statuspage.Render(
		w,
		r,
		http.StatusOK,
		statuspage.TemplateValues{
			Title:      "Signed in",
			Heading:    "Signed in",
			SubHeading: "You have signed in successfully",
			Paragraphs: []any{
				"You may close this web page.",
			},
		},
	)
}
