# ones-demo-2024
The code in this repo separates the [alvarium/example-go](https://github.com/project-alvarium/example-go) application into three
separate services for deployment and testing on different hosts. The Creator services instantiates sample data which is then
passed to the Mutator service via REST API, and then the Mutator service passes data to the Transitor service also via REST API.
This equates to the manner in which data is passed between goroutines in the example-go application, but allows for the capture
and annotation of TLS connectivity. By default, only the Transitor service has TLS enabled which allows us to demonstrate both
PASS/FAIL for the [TLS Annotator](https://github.com/project-alvarium/alvarium-sdk-go/blob/main/internal/annotators/tls.go). An
example of how this can be set up in a lab can be found on the first few pages of the [slide deck](https://github.com/project-alvarium/scoring-apps-go/blob/main/DCF_Scoring_Method_Proposed.pdf)
documenting the DCF Scoring Flow.

You can run all three of the services on a single box if you want to, but the hostname captured in the annotations will all be
the same which may or may not be suitable for your testing. If you only want to run a single service on a given host, then 
comment out the non-runnable services in the relevant [launch scripts](https://github.com/tsconn23/ones-demo-2024/tree/main/scripts/bin).
The launch scripts are triggered from the Makefile found in the root of this project.

Stream publication for the Alvarium annotations may be captured via either MQTT or IOTA Streeams. Every service's main() function
is hosted in the relevant folder under /cmd. From there, see the cmd/res directory for relevant config files that can be 
chosen when starting the service. 

More information for starting the services manually can be found in the launch scripts referenced above, as well as in the 
`go-example` application's [README](https://github.com/project-alvarium/example-go/blob/main/README.md).
