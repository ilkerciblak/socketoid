return {
  {
    "mason-org/mason.nvim",
    lazy = false,
    config = function()
      require("mason").setup()
    end,
  },

  {
    "mason-org/mason-lspconfig.nvim",
    dependencies = {
      "mason-org/mason.nvim",
      "neovim/nvim-lspconfig",
    },
    config = function()
      require("mason-lspconfig").setup({
        ensure_installed = {
          "lua_ls",
          "gopls",
          "ts_ls",
          "html",
          "emmet_ls",
        },
        automatic_installation = true,
      })
    end,
  },

  {
    "neovim/nvim-lspconfig",
    config = function()
      local capabilities = require("cmp_nvim_lsp").default_capabilities()

      -- LUA
      vim.lsp.config("lua_ls", {
        capabilities = capabilities,
        settings = {
          Lua = {
            diagnostics = {
              globals = { "vim" },
            },
          },
        },
      })


      -- GO
      vim.lsp.config("gopls", {
        capabilities = capabilities,
        settings = {
          gopls = {
            analyses = {
              unusedparams = true,
              shadow = true,
            },
            staticcheck = true,
            usePlaceholders = true,
          },
        },
      })

      -- TYPESCRIPT / TSX
      vim.lsp.config("ts_ls", {
        capabilities = capabilities,
      })

      -- HTML
      vim.lsp.config("html", {
        capabilities = capabilities,
      })

      -- EMMET (HTML snippet for TSX)
      vim.lsp.config("emmet_ls", {
        capabilities = capabilities,
        filetypes = {
          "html",
          "css",
          "scss",
          "javascriptreact",
          "typescriptreact",
        },
      })

      vim.lsp.enable("lua_ls")
      vim.lsp.enable("gopls")
      vim.lsp.enable("ts_ls")
      vim.lsp.enable("html")
      vim.lsp.enable("emmet_ls")


      -- Keymaps
      vim.keymap.set("n", "K", vim.lsp.buf.hover, {})
      vim.keymap.set("n", "<leader>gd", vim.lsp.buf.definition, {})
      vim.keymap.set({ "n", "v" }, "<leader>ca", vim.lsp.buf.code_action, {})
    end,
  },
}
