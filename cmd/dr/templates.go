package dr

const CascTemplateString = `
unclassified:
  location:
    url: http://server_ip:8080/
jenkins:
  securityRealm:
    local:
      allowsSignup: false
      users:
        - id: ${JENKINS_ADMIN_ID}
          password: ${JENKINS_ADMIN_PASSWORD}
  authorizationStrategy: unsecured
`

const RunScript = `
#!/usr/bin/env bash
set -euo pipefail
ITZ_OC_FLAGS="--insecure-skip-tls-verify=true"
# TODO: add some error handling here
oc login -u ${ITZ_OC_USER} -p ${ITZ_OC_PASS} ${ITZ_OC_URL} ${ITZ_OC_FLAGS} 
echo "Applying pipeline..."
oc apply -f ${ITZ_PIPELINE}
echo "Applying pipeline run..."
oc create -f ${ITZ_PIPELINE_RUN}

`
