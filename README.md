<picture align="center">
  <img src="https://user-images.githubusercontent.com/45081667/227194559-57c73a9e-51df-4d42-b9f0-8b6dd011ae7b.png" alt="Header" >
</picture>

<p align="center">
  <a href="https://1password.com">
    <h1 align="center">1Password Shell Plugins</h1>
  </a>
</p>

<p align="center">
 <h4 align="center"> Authenticate any CLI with your fingerprint </h4>
</p>

<p align="center">
  <a href="https://developer.1password.com/docs/cli/shell-plugins/"><img alt="Get Started" src="https://user-images.githubusercontent.com/45081667/226940040-16d3684b-60f4-4d95-adb2-5757a8f1bc15.png" height="37" /></a>
  <a href="https://developer.1password.com/docs/cli/shell-plugins/contribute"><img alt="Build" src="https://user-images.githubusercontent.com/45081667/226941089-473be407-5417-48e3-8710-55b2d4ec761a.png" height="37" /></a>
</p>
<br/>

---

Eliminate plaintext credentials in your home directory. Automatically authenticate every CLI you use with [1Password CLI](https://developer.1password.com/docs/cli/) + [Shell Plugins](https://developer.1password.com/docs/cli/shell-plugins/). Approve temporary credential usage in your terminal with biometrics.

## 🪄 Quick demo
<p align="center">
  <picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/45081667/227191964-9629476d-a49e-475d-b8cb-2115c302025d.gif">
   <img src="https://user-images.githubusercontent.com/45081667/227197994-6fdb2cad-c240-4cb7-ba4c-77f2483606ab.gif" alt="Example of 1Password Shell Plugins with AWS: user runs an `aws` command, a Touch ID prompt shows up, and `aws` is automatically authenticated" style="max-width: 100%; display: inline-block;" />
</picture>
</p>


## 🚀 Get started with the available plugins
<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://user-images.githubusercontent.com/45081667/226968760-70d3f6b0-a3eb-4c75-a674-6fd136d7149a.png">
  <img alt="Picture changing depending on light mode." src="https://user-images.githubusercontent.com/45081667/226969008-0a3f7537-7942-442f-9170-18b008a6574c.png">
</picture>
</p>

<a href="https://developer.1password.com/docs/cli/shell-plugins">
    <p align="center">See all...</p>
  </a>
<br/>

## 🔩 Don’t see yours? Contribute! <sup><b><a href="#-beta-notice">[BETA]</a></b></sup>
Is your favorite CLI not listed yet? Learn [how to build a new plugin](https://developer.1password.com/docs/cli/shell-plugins/contribute) yourself and then open a PR on this repository to get it added to 1Password!

For the contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

**Quick start:** clone this repo and run
```shell
make new-plugin
```

Still not sure where or how to get started? We're happy to help! You can:
- Book a free [pairing session](https://calendly.com/d/grs-x2h-pmb/1password-shell-plugins-pairing-session) with one of our developers
- Join the [Developer Slack workspace](https://join.slack.com/t/1password-devs/shared_invite/zt-1halo11ps-6o9pEv96xZ3LtX_VE0fJQA), and ask us any questions there

## 💙 Community & Support

- File an [issue](https://github.com/1Password/shell-plugins/issues/new/choose) for bugs and feature requests
- Join the [Developer Slack workspace](https://join.slack.com/t/1password-devs/shared_invite/zt-1halo11ps-6o9pEv96xZ3LtX_VE0fJQA)
- Subscribe to the [Developer Newsletter](https://1password.com/dev-subscribe/)

## 📣 Contributions Beta Notice

The 1Password Shell Plugins ecosystem is still in beta. In practice, this means that if you're building your own plugins locally, you'll likely have to recompile them on occasion to keep up with the latest updates of the 1Password CLI.