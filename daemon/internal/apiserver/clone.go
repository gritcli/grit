package apiserver

import (
	"github.com/gritcli/grit/api"
	"google.golang.org/protobuf/proto"
)

// CloneRepo makes a local clone of a repository from a source.
func (s *Server) CloneRepo(
	req *api.CloneRepoRequest,
	stream api.API_CloneRepoServer,
) error {
	repo, err := s.Cloner.Clone(
		stream.Context(),
		req.Source,
		req.RepoId,
		s.newClientLog(
			stream,
			req.ClientOptions,
			func(out *api.ClientOutput) proto.Message {
				return &api.CloneRepoResponse{
					Response: &api.CloneRepoResponse_Output{
						Output: out,
					},
				}
			},
		),
	)
	if err != nil {
		return err
	}

	return stream.Send(&api.CloneRepoResponse{
		Response: &api.CloneRepoResponse_LocalRepo{
			LocalRepo: marshalLocalRepo(repo),
		},
	})
}
