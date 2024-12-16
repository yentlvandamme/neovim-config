require('rose-pine').setup()

function SetupTheme(color)
    color = color or "rose-pine"
    vim.cmd.colorscheme(color)
end

SetupTheme()
