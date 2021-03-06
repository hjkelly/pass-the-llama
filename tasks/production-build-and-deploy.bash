rm -rf $GO_LINUX_BIN
mkdir $GO_LINUX_BIN
env GOOS=linux GOARCH=amd64 go build -o ${GO_LINUX_BIN}/sendPrompts ./commands/sendPrompts/
env GOOS=linux GOARCH=amd64 go build -o ${GO_LINUX_BIN}/processIncoming ./commands/processIncoming/
cat ${GOPATH}/../ptl-prod-inventory
ansible-playbook -i ${GOPATH}/../ptl-prod-inventory tasks/deploy.yml
