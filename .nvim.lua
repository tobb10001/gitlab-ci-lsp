local configs = require('lspconfig.configs')
local lspconfig = require('lspconfig')
local util = require('lspconfig.util')

configs.gitlab_ci_lsp = {
    default_config = {
        cmd = {"/usr/local/go/bin/go", "run", "main.go"},
        filetypes = {'gitlab-ci.yaml'},
        root_dir = util.find_git_ancestor,
        settings = {},
    },
}

lspconfig.gitlab_ci_lsp.setup({})
