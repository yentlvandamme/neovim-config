# neovim-config

Dear future me, this is a basic Neovim configuration based on the configuration of ThePrimeagen.
Keep in mind that keeping your Neovim config updated is not something which should be done continuously. You don't always need the latest of the latest. If something works, stick to it. You can take some time to upgrade all the plugins once every 6 months, unless you feel the need to do it earlier.
Make sure to keep your Neovim config as lightweight and minimalistic as possible. What I mean by that is that you should not download and configure plugins you'll never use. The power is Neovim is to configure only the stuff you need, and not have the stuff you don't need be in your face.

For future configurations, feel free to add new folders within this repository. It could be that you want a different configuration for side projects and for work.


## Requirements
There are certain technical requirements at the time of writing this:

- The Neovim version should at least be `0.9.0`. This version is not available through download. Meaning you'll have to download the source and build it from there. There are plenty of guides if needed
- We currently use Packer to install and sync all of the plugins. Make sure you have Packer installed.
