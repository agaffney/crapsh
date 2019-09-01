package flags

import ()

type Flag struct {
	Short       string
	Long        string
	Default     bool // whether flag is set by default
	CmdlineOnly bool // whether a flag is only available from the initial commandline
	HasArg      bool // whether the flag should be followed by an argument
}

var Flags = []Flag{
	{Long: `allexport`, Short: `a`},
	{Long: `braceexpand`, Short: `B`, Default: true},
	{Long: `emacs`, Short: ``, Default: true},
	{Long: `errexit`, Short: `e`},
	{Long: `errtrace`, Short: `E`},
	{Long: `functrace`, Short: `T`},
	{Long: `hashall`, Short: `h`, Default: true},
	{Long: `histexpand`, Short: `H`, Default: true},
	{Long: `history`, Short: ``, Default: true},
	{Long: `ignoreeof`, Short: ``},
	{Long: `interactive-comments`, Short: ``, Default: true},
	{Long: `keyword`, Short: `k`},
	{Long: `monitor`, Short: `m`, Default: true},
	{Long: `noclobber`, Short: `C`},
	{Long: `noexec`, Short: `n`},
	{Long: `noglob`, Short: `f`},
	{Long: `nolog`, Short: ``},
	{Long: `notify`, Short: `b`},
	{Long: `nounset`, Short: `u`},
	{Long: `onecmd`, Short: `t`},
	{Long: `physical`, Short: `P`},
	{Long: `pipefail`, Short: ``},
	{Long: `posix`, Short: ``},
	{Long: `privileged`, Short: `p`},
	{Long: `verbose`, Short: `v`},
	{Long: `vi`, Short: ``},
	{Long: `xtrace`, Short: `x`},
	{Short: `c`, CmdlineOnly: true, HasArg: true},
	{Short: `i`, CmdlineOnly: true},
	{Short: `l`, Long: `login`, CmdlineOnly: true},
	{Short: `r`, Long: `restricted`, CmdlineOnly: true},
	{Short: `s`, CmdlineOnly: true},
	{Short: `D`, Long: `dump-strings`, CmdlineOnly: true},
	{Long: `noprofile`, CmdlineOnly: true},
	{Long: `norc`, CmdlineOnly: true},
	{Long: `posix`, CmdlineOnly: true},
}
