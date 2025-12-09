package visualise

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/kvql/bunsceal/pkg/o11y"
)

// CheckGit checks if git binary is available
func CheckGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// HasGitHistory checks if the repository has full git history (not a shallow clone)
func HasGitHistory() (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--is-shallow-repository")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	// "false" means NOT shallow, i.e., has full history
	return string(output) == "false\n", nil
}

func GetLatestCommitTime(path string) (time.Time, error) {
	var pathT time.Time
	epochRe := regexp.MustCompile("[0-9]+")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		message := fmt.Sprintf("File/dir does not exist, %s", path)
		cp, _ := os.Getwd()
		o11y.Log.Printf("current path: %s", cp)
		return pathT, errors.New(message)
	}
	txOut := exec.Command("git", "--no-pager", "log", "-1", "--format=\"%cd\"", "--date=raw", "--", path)
	tmpOut, err := txOut.Output()
	if err != nil {
		o11y.Log.Println("Error checking latest taxonomy commit date. Error: ", err)
		o11y.Log.Println(string(tmpOut))
		return pathT, errors.New("error checking latest taxonomy commit date")
	}
	// sample format "1710150183 +0000", below is ignoring the timezone
	tmp := epochRe.FindStringSubmatch(string(tmpOut[:]))
	if len(tmp) == 0 {
		o11y.Log.Println("Error parsing time")
		return pathT, errors.New("error parsing time, no match")
	}
	// convert to int
	tmpInt, err := strconv.Atoi(tmp[0])
	if err != nil {
		return pathT, errors.New("error converting time to int")
	}
	pathT = time.Unix(int64(tmpInt), 0)
	return pathT, nil
}
