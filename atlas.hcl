data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./models",
    "--dialect", "postgres"
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgresql://micro:dranreb22@localhost:5432/goauth?search_path=public&sslmode=disable"
  url = "postgresql://micro:dranreb22@localhost:5432/goauth?search_path=public&sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}