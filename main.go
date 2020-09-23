package main

import (
	"log"
	"strconv"

	"github.com/codersrank-org/multi_repo_repo_extractor/flag"
	repo "github.com/codersrank-org/multi_repo_repo_extractor/repo"
	upload "github.com/codersrank-org/multi_repo_repo_extractor/upload"
)

func main() {

	// TODO implement auto-update (versioning etc.)

	provider, token, repoVisibility, repoInfoExtractorPath, emails := flag.ParseFlags()

	repositoryService := repo.NewRepositoryService(repoInfoExtractorPath, provider, repoVisibility, token, emails)
	codersrankService := upload.NewCodersrankService()

	repositoryService.InitRepoInfoExtractor()
	repos := repositoryService.GetReposFromProvider()
	uploadResults := make(map[string]string)

	for _, repo := range repos {
		err := repositoryService.Clone(repo)
		if err != nil {
			log.Printf("Couldn't clone repo, skipping: %s, error: %s", repo.FullName, err.Error())
			continue
		}
		err = repositoryService.Process(repo)
		if err != nil {
			log.Printf("Couldn't process repo, skipping: %s, error: %s", repo.FullName, err.Error())
			continue
		}
		uploadToken, err := codersrankService.UploadRepo(strconv.Itoa(repo.ID))
		if err != nil {
			log.Printf("Couldn't upload processed repo: %s, error: %s", repo.FullName, err.Error())
			continue
		}
		uploadResults[repo.Name] = uploadToken
	}

	resultToken := codersrankService.UploadResults(uploadResults)
	codersrankService.ProcessResults(resultToken)

}