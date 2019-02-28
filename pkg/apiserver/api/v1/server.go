package v1

import (
	"context"
	errs "errors"
	"fmt"
	"log"

	"minibox.ai/pkg/errors"

	pb "minibox.ai/pkg/api/v1/core"
	"minibox.ai/pkg/database"
	"minibox.ai/pkg/models"
	"minibox.ai/pkg/utils"
)

type Server struct{}

var db *database.Database

func SetDB(_db *database.Database) {
	db = _db
}

func panicdb() {
	if db == nil {
		panic("you must call SetDB(db* database.Database), db can't be nil")
	}
}

// SayHello implements api.ApiServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// SayHello implements api.ApiServer
func (s *Server) ListProjects(ctx context.Context, in *pb.ProjectsRequest) (*pb.ProjectsReply, error) {
	panicdb()
	var (
		prjs    []models.Project
		cur_usr *models.User
		ok      bool
	)

	if cur_usr, ok = getCurrntUser(ctx); !ok {
		return nil, errs.New("not have available User, please use `mini login` as a authorize User")
	}

	if err := db.Find(&prjs, "user_id = ?", cur_usr.ID).Error; err != nil {
		return nil, err
	}

	return &pb.ProjectsReply{Items: convertProjects(prjs)}, nil
}

func (s *Server) CreateProject(ctx context.Context, in *pb.CreateProjectRequest) (*pb.Project, error) {
	var (
		cur_usr *models.User
		// uid     = uuid.NewV4()
		ok bool
	)
	panicdb()

	if cur_usr, ok = getCurrntUser(ctx); !ok {
		return nil, errs.New("not have available User, please use `mini login` as a authorize User")
	}

	var prj = models.Project{
		Model:     models.Model{ID: utils.GenerateRandomID()},
		Name:      in.Name,
		Namespace: cur_usr.Namespace,
		UserID:    cur_usr.ID,
	}

	if !db.First(&models.Project{}, "name = ? and namespace = ?", in.Name, cur_usr.Namespace).RecordNotFound() {
		return nil, fmt.Errorf("the '%s' named is exsits, use another name again.", in.Name)
	}

	if err := db.Create(&prj).Error; err != nil {
		log.Printf("create Project error %s", err)
		return nil, err
	}

	return &pb.Project{Id: prj.ID, Name: prj.Name}, nil
}

func (s *Server) UpdateProject(ctx context.Context, in *pb.UpdateProjectRequest) (*pb.UpdateProjectReply, error) {
	return nil, errors.ErrNonImplement
}

func (s *Server) RemoveProject(ctx context.Context, in *pb.RemoveProjectRequest) (*pb.RemoveProjectReply, error) {
	return nil, errors.ErrNonImplement
}

// project train
func (s *Server) TrainStart(ctx context.Context, in *pb.TrainStartRequest) (*pb.TrainReply, error) {
	return nil, errors.ErrNonImplement

}
func (s *Server) TrainPause(ctx context.Context, in *pb.TrainPauseRequest) (*pb.TrainReply, error) {
	return nil, errors.ErrNonImplement

}
func (s *Server) TrainResume(ctx context.Context, in *pb.TrainResumeRequest) (*pb.TrainReply, error) {
	return nil, errors.ErrNonImplement

}
func (s *Server) TrainStop(ctx context.Context, in *pb.TrainStopRequest) (*pb.TrainReply, error) {
	return nil, errors.ErrNonImplement

}

// func (s *Server) data management
func (s *Server) DataUploadFile(ctx context.Context, in *pb.DataUploadRequest) (*pb.UploadReply, error) {
	return nil, errors.ErrNonImplement

}
func (s *Server) DataUploadFromS3(ctx context.Context, in *pb.DataUploadS3Request) (*pb.UploadReply, error) {
	return nil, errors.ErrNonImplement

}
func (s *Server) DataDownload(ctx context.Context, in *pb.DataDownloadRequest) (*pb.DownloadReply, error) {
	return nil, errors.ErrNonImplement

}
func (s *Server) DataScrapingFromURL(ctx context.Context, in *pb.DataScrapingRequest) (*pb.ScrapingReply, error) {
	return nil, errors.ErrNonImplement

}

// func (s *Server) pod settings path
func (s *Server) SetupMount(ctx context.Context, in *pb.MountRequest) (*pb.MountReply, error) {
	return nil, errors.ErrNonImplement
}

// func (s *Server) hyper parameters set
func (s *Server) SetParameterLearningRate(ctx context.Context, in *pb.HyperParameterRequest) (*pb.HyperParameterReply, error) {
	return nil, errors.ErrNonImplement
}

// SayHello implements api.ApiServer
func (s *Server) CreateDataset(ctx context.Context, in *pb.DatasetName) (*pb.DatasetReply, error) {
	panicdb()
	var (
		cur_usr *models.User
		ok      bool
		err     error
	)

	if cur_usr, ok = getCurrntUser(ctx); !ok {
		return nil, errs.New("not have available User, please use `mini login` as a authorize User")
	}

	var dataset = models.Dataset{
		Namespace: cur_usr.Namespace,
		Name:      in.Name,
	}

	log.Printf("namespace %s and name %s", cur_usr.Namespace, in.Name)
	if db.First(&models.Dataset{}, "namespace = ? and name = ?", cur_usr.Namespace, in.Name).RecordNotFound() {

		if err = db.Create(&dataset).Error; err != nil {
			return nil, err
		}
		return &pb.DatasetReply{Name: dataset.Name}, nil
	} else {
		return nil, errors.ErrAlreadyRecord
	}
}

func (s *Server) ListDatasets(ctx context.Context, in *pb.ListDatasetsRequest) (*pb.DatasetsReply, error) {
	panicdb()
	var (
		cur_usr  *models.User
		ok       bool
		datasets []models.Dataset
		err      error
	)

	if cur_usr, ok = getCurrntUser(ctx); !ok {
		return nil, errs.New("not have available User, please use `mini login` as a authorize User")
	}

	log.Printf("ListDatasets namespace %s", cur_usr.Namespace)

	if err = db.Find(&datasets, "namespace = ? ", cur_usr.Namespace).Error; err != nil {
		return nil, err
	}

	var reply pb.DatasetsReply
	for _, ds := range datasets {
		name := &pb.DatasetName{Name: ds.Name}
		reply.Datasets = append(reply.Datasets, name)
	}

	return &reply, nil

}

func (s *Server) GetDataset(ctx context.Context, in *pb.GetDatasetRequest) (*pb.DatasetReply, error) {
	return nil, errors.ErrNonImplement
}

func (s *Server) RegistryDataset(ctx context.Context, in *pb.RegistryDatasetRequest) (*pb.DatasetReply, error) {
	return nil, errors.ErrNonImplement
}

func (s *Server) GetFileS3Object(ctx context.Context, in *pb.GetFileS3Request) (*pb.FileObjectReply, error) {
	return nil, errors.ErrNonImplement
}

func (s *Server) ListDatasetObjects(ctx context.Context, in *pb.ListDatasetFilesRequest) (*pb.FilesReply, error) {
	return nil, errors.ErrNonImplement
}

func convertProjects(prjs []models.Project) []*pb.Project {
	var pbprjs = make([]*pb.Project, len(prjs))
	for i, prj := range prjs {
		pbprjs[i] = &pb.Project{
			Id:      prj.ID,
			Name:    prj.Name,
			Summary: prj.Summary,
		}
	}

	return pbprjs
}
