# Contributing

Banji is, all things considered, a boring project. It is not a software product; it is a library. A library that just
aims to provide a framework for the way data is moved around, in the most reductive of terms.

Therefore, the fact that you have taken an interest in this project is quite cool! I am always looking for ways to make
Banji more useful, performant, and sensible. Another set of eyes on a project is always beneficial.

I do not have extensive experience in maintaining open-source projects, so please keep that in mind.

## Test Coverage

Banji is designed for use in both consecutive and concurrent environments; therefore, the tests we write for features
and changes should reflect this reality. We should strive for complete coverage, not only to ensure Banji functions as
intended, but also to benchmark its performance.

Contributing to Banji's test coverage is incredibly valuable, just as much as fixing bugs and implementing new features.

## Code Conventions

- Indents should be four spaces, not eight or two.
- Avoid nesting whenever possible. Nested code is harder to read and maintain.
- Prioritize clarity over brevity.
- External dependencies should, in most cases, be avoided.
- When instantiating structs, fields should be explicit, not implicit.
- Magic strings and numbers should be extracted to constants or variables.

## Submitting Changes

Banji uses a fork-and-pull model. Please send pull requests to @AndrewChon. Ensure that pull requests detail all the
changes you have made and the impetus for those changes. Commits should have descriptive messages that aptly summarize
what changes have been made.

