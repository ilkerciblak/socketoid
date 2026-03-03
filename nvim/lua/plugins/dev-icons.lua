return {
	"nvim-tree/nvim-web-devicons",
	config = function()
		require("nvim-web-devicons").setup({
			default = true,
		})
		require("nvim-web-devicons").get_icons()
	end,
}
