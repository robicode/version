# Version: Go Port of Ruby Gem::Version

This is a mostly complete port of rubygems' `Gem::Version` and `Gem::Requirement` classes to Go constructs, with concessions made for Go coding standards.

The library tries to stay true to the Ruby implementation. Please see the source for the interface, and the TODO.md file for what needs to be done.

## Basic Usage

For the `*Version` construct:

``` go
// Initialize a new *Version
ver, err := version.New("1.5.3")

// If you know the version string is valid and don't care about err
// (use with caution!). ver will be nil on error:
ver := New2("1.5.3")

// Return the parsed version as a string
myString := ver.Version()

// Compare with another *Version
ver.Compare(ver2)

// Check if ver is a prerelease
isPre := ver.IsPrerelease()

// Get an approximate recommendation for a tilde requirement:
v := ver.ApproximateRecommendation()

// Bump a version
newVersion := ver.Bump()
```

See the source for more functionality.

## Requirement Construct

``` go
// Create a new *Requirement construct
requirement, err := requirement.New(">= 1.3", "< 1.4")

// Checks is a *Version construct meets all supplied requirements
result := requirement.IsSatisfiedBy(ver)
```

## Contributing

Just submit a PR or an issue.

## License

Licensed under MIT as per the LICENSE file.
