echo "Initializing Neovim config"

config_target=~/.config/nvim
source=./distributions/Azalea/.

echo "Target location of config:" $config_target
echo "Source of config:" $source

rm -r $config_target
mkdir $config_target
cp -R $source $config_target
