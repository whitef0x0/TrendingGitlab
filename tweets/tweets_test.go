package tweets

import (
	"net/url"
	"testing"

	"github.com/whitef0x0/TrendingGitlab/github"
	"github.com/whitef0x0/go-trending"
)

func TestTweets_IsProjectEmpty(t *testing.T) {
	ts := TweetSearch{}
	mock := []struct {
		Project trending.Project
		Result  bool
	}{
		{trending.Project{Name: ""}, true},
		{trending.Project{Name: "MyProject"}, false},
	}

	for _, item := range mock {
		res := ts.IsProjectEmpty(item.Project)
		if res != item.Result {
			t.Errorf("Failed for project \"%s\", got %v, expected %v", item.Project.Name, res, item.Result)
		}
	}
}

func TestTweets_BuildTweet(t *testing.T) {
	owner := "whitef0x0"
	repositoryName := "TrendingGitlab"
	projectName := owner + "/" + repositoryName
	projectURL, _ := url.Parse("https://github.com/whitef0x0/TrendingGitlab")
	projectDescription := "A twitter bot (@TrendingGitlab) to tweet trending repositories and developers from GitHub"

	ts := TweetSearch{
		URLLength: 24,
	}

	stars := 123
	repository := &github.Project{
		StarCount: stars,
	}

	mock := []struct {
		Project trending.Project
		Result  string
	}{
		{trending.Project{
			Stars: stars,
			Name:           "SuperDuperOwnerOrOrganisation/This-Is-A-Long-Project-Name-That-Will-Drop-The-Description-Of-The-Project",
			Owner:          "SuperDuperOwnerOrOrganisation",
			RepositoryName: "This-Is-A-Long-Project-Name-That-Will-Drop-The-Description-Of-The-Project",
			NameSpace: 		"superduper/long-project-name",
			Description:    projectDescription + " and more and better and super duper text",
			URL:            projectURL,
		}, "SuperDuperOwnerOrOrganisation/This-Is-A-Long-Project-Name-That-Will-Drop-The-Description-Of-The-Project ★123 https://github.com/whitef0x0/TrendingGitlab"},
		{trending.Project{
			Stars: stars,
			Name:           projectName + "-cool-super-project",
			Owner:          owner,
			RepositoryName: repositoryName + "-cool-super-project",
			NameSpace: 		"whitef0x0/repo-cool-super-project-name",
			Description:    projectDescription + " and more and better and super duper text",
			URL:            projectURL,
		}, "whitef0x0/TrendingGitlab-cool-super-project: A twitter bot (@TrendingGitlab) to tweet trending... ★123 https://github.com/whitef0x0/TrendingGitlab"},
		{trending.Project{
			Stars: stars,
			Name:           projectName,
			Owner:          owner,
			RepositoryName: repositoryName,
			NameSpace: 		owner+"/"+repositoryName,
			Description:    projectDescription,
			URL:            projectURL,
		}, "whitef0x0/TrendingGitlab: A twitter bot (@TrendingGitlab) to tweet trending repositories and developers... ★123 https://github.com/whitef0x0/TrendingGitlab"},
		{trending.Project{
			Stars: stars,
			Name:           projectName,
			Owner:          owner,
			RepositoryName: repositoryName,
			NameSpace: 		owner+"/"+repositoryName,
			Description:    "Short description",
			URL:            projectURL,
		}, "whitef0x0/TrendingGitlab: Short description ★123 https://github.com/whitef0x0/TrendingGitlab"},
		{trending.Project{
			Stars: stars,
			Name:           projectName,
			Owner:          owner,
			RepositoryName: repositoryName,
			Description:    "Project without a URL",
		}, "whitef0x0/TrendingGitlab: Project without a URL ★123"},
		{trending.Project{
			Stars: stars,
			Name:           repositoryName + "/" + repositoryName,
			Owner:          repositoryName,
			RepositoryName: repositoryName,
			NameSpace: 		repositoryName+"/"+repositoryName,
			Description:    projectDescription,
			URL:            projectURL,
		}, "TrendingGitlab: A twitter bot (@TrendingGitlab) to tweet trending repositories and developers from GitHub ★123 https://github.com/whitef0x0/TrendingGitlab"},
	}

	for _, item := range mock {
		res := ts.BuildTweet(item.Project, repository)
		if res != item.Result {
			t.Errorf("Failed building a tweet for project \"%s\". Got \"%s\", expected \"%s\"", item.Project.Name, res, item.Result)
		}
	}
}

var testSlice = []string{"one", "two", "three", "four"}

func TestUtility_ShuffleStringSlice_Length(t *testing.T) {
	shuffledSlice := make([]string, len(testSlice))
	copy(shuffledSlice, testSlice)
	ShuffleStringSlice(shuffledSlice)

	if len(testSlice) != len(shuffledSlice) {
		t.Errorf("The length of slices are not equal. Got %d, expected %d", len(shuffledSlice), len(testSlice))
	}
}

func TestUtility_ShuffleStringSlice_Items(t *testing.T) {
	shuffledSlice := make([]string, len(testSlice))
	copy(shuffledSlice, testSlice)
	ShuffleStringSlice(shuffledSlice)

	for _, item := range testSlice {
		if IsStringInSlice(item, shuffledSlice) == false {
			t.Errorf("Item \"%s\" not found in shuffledSlice: %+v", item, shuffledSlice)
		}
	}
}

func TestUtility_Crop(t *testing.T) {
	testSentence := "This is a test sentence for the unit test."
	textMock := []struct {
		Content     string
		Chars       int
		AfterString string
		Crop2Space  bool
		Result      string
	}{
		{testSentence, 0, "", false, testSentence},
		{testSentence, 99, "", false, testSentence},
		{testSentence, 13, "", false, "This is a te"},
		{testSentence, 13, "...", false, "This is a te..."},
		{testSentence, 13, "", true, "This is a"},
		{testSentence, 13, "...", true, "This is a..."},
		{testSentence, -99, "", false, testSentence},
		{testSentence, -13, "", false, "he unit test."},
		{testSentence, -13, "...", false, "...he unit test."},
		{testSentence, -13, "", true, "unit test."},
		{testSentence, -13, "...", true, "...unit test."},
	}

	for _, mock := range textMock {
		res := Crop(mock.Content, mock.Chars, mock.AfterString, mock.Crop2Space)
		if res != mock.Result {
			t.Errorf("Crop result is \"%s\", but expected \"%s\".", res, mock.Result)
		}
	}
}

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
