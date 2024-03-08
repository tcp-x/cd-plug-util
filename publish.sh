# get current repository latest version
echo "current repository latest version:\n"
git ls-remote --tags https://github.com/tcp-x/cd-plug-util.git
# set latest version
Version="v0.0.8"

# cd $projDir
go mod tidy
git add go.mod go.sum i-exec.go
git commit -am "set version $Version"
git tag $Version
git push origin $Version

