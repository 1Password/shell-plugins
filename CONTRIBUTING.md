# 👫 Contributing to 1Password Shell Plugins

Thanks for your interest in contributing to the 1Password Shell Plugins project! 🙌

The only way we can make plaintext secrets in the home directory a thing of the past and get the widest coverage of shell plugins is by making this a community effort.

This document contains useful information and guidelines for everyone interested in contributing to the project.

<!----><a name="your-first-contribution"></a>
## 🎉 Your First Plugin

If you're planning to create your first plugin, we recommend you have a look at the [step-by-step guide](https://developer.1password.com/docs/cli/shell-plugins/contribute/) in the 1Password developer docs first, which explains the basic concepts of the plugin ecosystem and walks you through the creation of your first plugin.

<!----><a name="scope"></a>
## 🔭 Scope

The current focus of 1Password Shell Plugins is to cover CLIs that require authentication to their platform's backend.
Examples of that are CLIs provided by SaaS platforms, cloud vendors, and databases.

Unofficial CLIs are less likely to get accepted at this moment.

CLIs that authenticate to multiple different platforms depending on the project are also out of scope at this stage.
Examples include `terraform`, `ansible`, or the CLIs that come with application development frameworks (`flask run`, `spring run`, etc.).

<!----><a name="third-party-dependencies"></a>
## 🖇️ Third-Party Dependencies

Try to avoid bringing in third-party dependencies if possible.
Especially when building importers for your plugin, all that's needed in most cases is to read a file on disk, for which the shell plugins SDK already provides [helpers](sdk/importer/).

PRs that introduce new dependencies will take longer to get through the review process and have a higher chance of getting rejected.

<!----><a name="sign-your-commits"></a>
## ✍️ Sign Your Commits

To get your PR merged, we require you to sign your commits.
Fortunately, this has become [very easy to set up](https://developer.1password.com/docs/ssh/git-commit-signing)!

<!----><a name="testing"></a>
## 🧪 Testing

The plugin ecosystem allows you to [locally build](#make-plugin-build) your plugins and use them with 1Password CLI, so you can take your plugin for a spin before submitting a PR.

You can add tests to your plugin using the SDK's [`plugintest` package](sdk/plugintest/), which provides helpers so that you only have to care about the test cases themselves.
You can use the [`example-secrets` command](#make-plugin-example-secrets) to help create test fixtures.

<!----><a name="makefile-commands"></a>
## 👷 Makefile Commands

This repo comes with a set of `make` commands to make your life as a contributor easier.

<!----><a name="make-new-plugin"></a>
### Create a New Plugin

Bootstrap a new plugin package based on a few prompts to fill in:

```
make new-plugin
```

<!----><a name="make-plugin-validate"></a>
### Validate Plugin Schema

Validate the schema of a plugin and prints out a report of all validation checks that were done:

```
make <plugin>/validate
```

<!----><a name="make-plugin-build"></a>
### Locally Build Your Plugin

Locally build a plugin and install it in the `~/.op/plugins/local` directory:

```
make <plugin>/build
```

**Note:** Locally built plugins take precedence over existing plugins. This allows you to customize existing plugins.

<!----><a name="make-plugin-example-secrets"></a>
### Print Example Secrets

Automatically generate example values for all secrets defined in a plugin:

```
make <plugin>/example-secrets
```

<!----><a name="get-in-touch"></a>

## 📄 Documentation

No documentation work is required on your end. If your plugin is accepted, a technical writer on our team will publish [a guide to using your plugin](https://developer.1password.com/docs/cli/shell-plugins) when the plugin is added to 1Password CLI. You can help the technical writing team by:
- Making note in your PR of anything unique to your shell plugin that might need to be documented beyond our standard template.
- Including an example command that requires authentication for us to use in our testing and screenshots.

## 💬 Get In Touch

If you need help, found a bug, or have an idea for an awesome plugin that you'd like to discuss with us first, you can reach out to us here:

* Create an [issue](https://github.com/1Password/shell-plugins/issues) in this repo.
* Join the [Developer Slack workspace](https://developer.1password.com/joinslack).

## 📣 Contributions Beta Notice

The 1Password Shell Plugins ecosystem is still in beta. In practice, this means that if you're building your own plugins locally, you'll likely have to recompile them on occasion to keep up with the latest updates of the 1Password CLI.
