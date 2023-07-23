# Contributing to `docker-typo3`

Thank you for considering contributing to `docker-typo3`! Before you get started, please take a moment to review the guidelines for contributing to this repository. Since this project contains Dockerfiles that are auto-generated, it is important to follow specific procedures to maintain consistency and avoid conflicts.

## How Can I Contribute?

There are several ways to contribute to this set of container images. Whether you find a bug, have a feature request, or want to submit code changes, your help is valuable.

### Report Bugs

If you encounter a bug while using this container image, please [open an issue](../../issues) on this repository. Be sure to provide as much detail as possible, including steps to reproduce the bug and the expected behavior.

### Request Features

You can also use the [issue tracker](../../issues) to request new features or enhancements. Please check if a similar request already exists before creating a new one. Describe the feature request and why it would be valuable to the project.

### Contribute Code

Contributions are welcome! If you want to add a feature, fix a bug, or make any changes to the Dockerfiles, follow these steps:

1. **Do Not Edit Dockerfiles Directly**: Dockerfiles in this repository are auto-generated from a template.
2. **Edit the Template**: Make necessary changes to the template located in the [`updater`](updater/) directory.
3. **Generate Dockerfiles**: After modifying the template, run `make all` to produce updated Dockerfiles.
4. **Test Locally**: Test the generated Dockerfiles locally to ensure they work as expected.
5. **Submit a Pull Request**: Finally, submit a pull request with your changes.

## Generating Dockerfiles

The Dockerfiles in this repository are automatically generated a template. To generate Dockerfiles, follow these steps:

1. Make sure you have the Go SDK installed on your machine in order to run the code generatro
2. Make the required changes to the `Dockerfile.tmpl` located in the [`updater`](updater/) directory.
3. Run the generator by calling `make all`, which will take care of producing the updated Dockerfiles based on the templates.

Please note that manually editing the Dockerfiles in this repository will not be accepted. Always make changes to the template and follow the generation process.

## Submitting a Pull Request
When you're ready to submit your contributions, please follow these guidelines to ensure a smooth review process:

1. Fork this repository and create a new branch for your changes.
2. Make your changes following the guidelines mentioned above.
3. Test the generated Dockerfiles locally to verify that they are working as expected.
4. Commit your changes and provide a clear and descriptive commit message.
5. Push your branch to your forked repository.
6. Open a pull request (PR) against the main repository's `master` branch.

## License
By contributing to this project, you agree that your contributions will be licensed under the [project's license](LICENSE).

---

Thank you for your interest in contributing to the TYPO3 Docker image. Your efforts are greatly appreciated!

_Courtesy to ChatGTP for writing this contribution guideline. Glory to our future AI overlords. 🤖_