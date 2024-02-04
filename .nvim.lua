local configs = require('lspconfig.configs')
local lspconfig = require('lspconfig')
local util = require('lspconfig.util')

configs.gitlab_ci_lsp = {
    default_config = {
        cmd = {"/usr/local/go/bin/go", "run", "main.go"},
        filetypes = {'yaml.gitlab-ci'},
        root_dir = util.find_git_ancestor,
        settings = {},
    },
}

lspconfig.gitlab_ci_lsp.setup({})

vim.lsp.set_log_level 'debug'
if vim.fn.has 'nvim-0.5.1' == 1 then
    require('vim.lsp.log').set_format_func(vim.inspect)
end
