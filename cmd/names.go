package cmd

// The names and actions here, as constants, are to provide a consistent terminology for the command
// and action names.

const CreateAction = "create"
const DeployAction = "deploy"
const DoctorAction = "doctor"
const ExecuteAction = "execute"
const ListAction = "list"
const LoginAction = "login"
const ReserveAction = "reserve"
const ShowAction = "show"
const VersionAction = "version"

const ApiResource = "api"
const BuildResource = "build"
const EnvironmentResource = "environment"
const PipelineResource = "pipeline"
const ReservationResource = "reservation"
const WorkspaceResource = "workspace"

const TechZoneFull = "IBM Technology Zone"
const TechZoneShort = "TechZone"

const PipelineAnnotation = "techzone.ibm.com/tekton-pipeline-location"
const PipelineRunAnnotation = "techzone.ibm.com/tekton-pipeline-run-location"
