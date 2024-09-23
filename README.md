# Temporal.io Proof of Concept
This repository contains a system that uses Temporal.io as the workflow executor. The workflow execution is reentrant,
meaning it is resumable, recoverable, and reactive.


## How to Run
### Generate Protobuf
```shell
$ make generate package=api
```

### Run on Docker
Run below procedures in order to run the whole system completely
#### Run Service Dependencies
```shell
$ docker compose up
$ make migrate-up
```
#### Run Services
```shell
$ make build-image-all
$ docker compose -f app-docker-compose.yml up
```

### Run on Kubernetes
#### Run Services
```shell
$ make build-image-all
$ make prepare-kube
$ make deploy-kube-dependency
$ make deploy-kube-all
```

## How to Operate
### Monitor Workflow
#### Temporal.io
We use Temporal.io admin UI to monitor and interact with the workflow execution.
Open http://localhost:8080 in the browser.

### Hit gRPC
#### grpcui
We use grpcui to interace with the gRPC server provided by the transaction server.
```shell
$ grpcui -plaintext -port=9091 localhost:9090
```
Open http://localhost:9091 in browser

## System Flow
[![](https://mermaid.ink/img/pako:eNqtVsuOmzAU_ZUrrzPtnsVsUmmkqo-0pDukyuAbYgns1L7ONBrNv9dAEl6GoaVIRInle865j-P4hWVaIIuYxV8OVYYfJM8NLxMF_vlh0cDD4yM8ccJnfoEvfvPHOIItLwrYfY338J4MV5ZnJLVqgvp76_B9uwdiNGc0V4j8-24LW4M-ZD_EGQfVWDttKUcTf_sUQY4EBbcEHRGQXno_pWjgeEH-9VTiAvhbWrLNevW0mJNyDZIzqocsOPEWY0LusHJXHDRGm6AeLCyC0gQH7ZT4F40NdgBirqIG64pmdS9g1NS7gnrfpALrsgztLZOKe9zLrfEfQ96qk06dXFpIe0QxUecZBRXqG10KinjSOi8Qdi61Lo3gqmCEAKTBeTf89O9ZZvjQ2fCuKdo92do1DWCIItPKuhLHFAejywUkIYpuKQUKl1ENBCkvuDd164B6NiaHqgPbn6Y2Yo75OjxdZsAzKoJnSUc4cFkMWmuJk7PL9PRnq3JJb-WvwldmcvKjLVX-n1JR_clZ5I-xsqBD7oD3dgZSah3S55-wRgDBm6NTiuD4GqxWMXgivO2VCdWNZZZTzxAPi73gH2XmNBp4KHweVX5c4ZUFiYxP8wBHcHDn0xpbcY1TVidyI1mdiRJwm0-2YSWakkvhL0cv1VrC6IglJizyXwUeuCsoYYl69Vu5Ix1fVMYiMg43zJ18m293KRYduK_RhqGQpM3n5sJV37te_wB1ZEY1?type=png)](https://mermaid.live/edit#pako:eNqtVsuOmzAU_ZUrrzPtnsVsUmmkqo-0pDukyuAbYgns1L7ONBrNv9dAEl6GoaVIRInle865j-P4hWVaIIuYxV8OVYYfJM8NLxMF_vlh0cDD4yM8ccJnfoEvfvPHOIItLwrYfY338J4MV5ZnJLVqgvp76_B9uwdiNGc0V4j8-24LW4M-ZD_EGQfVWDttKUcTf_sUQY4EBbcEHRGQXno_pWjgeEH-9VTiAvhbWrLNevW0mJNyDZIzqocsOPEWY0LusHJXHDRGm6AeLCyC0gQH7ZT4F40NdgBirqIG64pmdS9g1NS7gnrfpALrsgztLZOKe9zLrfEfQ96qk06dXFpIe0QxUecZBRXqG10KinjSOi8Qdi61Lo3gqmCEAKTBeTf89O9ZZvjQ2fCuKdo92do1DWCIItPKuhLHFAejywUkIYpuKQUKl1ENBCkvuDd164B6NiaHqgPbn6Y2Yo75OjxdZsAzKoJnSUc4cFkMWmuJk7PL9PRnq3JJb-WvwldmcvKjLVX-n1JR_clZ5I-xsqBD7oD3dgZSah3S55-wRgDBm6NTiuD4GqxWMXgivO2VCdWNZZZTzxAPi73gH2XmNBp4KHweVX5c4ZUFiYxP8wBHcHDn0xpbcY1TVidyI1mdiRJwm0-2YSWakkvhL0cv1VrC6IglJizyXwUeuCsoYYl69Vu5Ix1fVMYiMg43zJ18m293KRYduK_RhqGQpM3n5sJV37te_wB1ZEY1)
