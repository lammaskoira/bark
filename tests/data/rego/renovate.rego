package bark

default allow := false

allow {
    file.exists("./renovate.json")
}