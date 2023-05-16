package kindsys

// Provider is a structure that allows packaging core, composable and
// custom kinds together with some metadata that uniqely identifies the
// provider.
Provider: {
  // The unique name of the provider.
  // TODO: What should the constraints be?
  name: =~"^([a-z][a-z0-9-]+?)$"

  // The version of the provider. Must be formatted according to semantic versioning
  // rules (<major>.<minor>.<patch>), e.g. 1.0.0.
  // TODO: What should the constraints be?
  version: =~"^([0-9]\\.[0-9]\\.[0-9])$"

  // Core kinds provided by the provider.
  coreKinds?: [Name=string]: Core & {
    name: Name
  }

  // Composable kinds provided by the provider.
  composableKinds?: [Iface=string]: Composable & {
    schemaInterface: Iface
  }

  // Custom kinds provided by the provider.
	customKinds?: [Name=string]: Custom & {
    name: Name
  }

  ...
}
