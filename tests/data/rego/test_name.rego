package bark

default allow := false

allow {
    contentstr := file.read(".github/workflows/test.yml")
    contents := yaml.unmarshal(contentstr)
    contents.name == "test"
}