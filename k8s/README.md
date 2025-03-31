If you want to install spotahome/redis-operator you may encounter an error:
```
Error: INSTALLATION FAILED: failed to install CRD crds/databases.spotahome.com_redisfailovers.yaml: error parsing : error converting YAML to JSON: yaml: line 4: did not find expected node content
```
For me [this comment](https://github.com/spotahome/redis-operator/issues/679#issuecomment-1853390076) under the github issue solved the problem.