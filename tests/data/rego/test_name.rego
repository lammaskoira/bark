package bark

default allow := false

allow {
    contentstr := file.readall(".github/workflows/test.yml")
    contents := yaml.unmarshal(contentstr)
    contents.name == "test"
}