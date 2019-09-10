# Estafette CI

The `estafette-extension-gcs` component is part of the Estafette CI system documented at https://estafette.io.

Please file any issues related to Estafette CI at https://github.com/estafette/estafette-ci-central/issues

## Estafette-extension-gcs

This extension helps to copy files into a Google Cloud Storage bucket.

## Development

To start development run

```bash
git clone git@github.com:estafette/estafette-extension-gcs.git
cd estafette-extension-gcs
```

Before committing your changes run

```bash
go test ./...
go mod tidy
```