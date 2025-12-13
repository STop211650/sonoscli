package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

type groupAudioClient interface {
	GetGroupVolume(ctx context.Context) (int, error)
	SetGroupVolume(ctx context.Context, volume int) error
	GetGroupMute(ctx context.Context) (bool, error)
	SetGroupMute(ctx context.Context, mute bool) error
}

var newGroupAudioClient = func(ctx context.Context, flags *rootFlags) (groupAudioClient, error) {
	return coordinatorClient(ctx, flags)
}

func newGroupVolumeCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Get or set group volume",
		Long:  "Controls GroupRenderingControl group volume on the group coordinator (0-100).",
	}

	cmd.AddCommand(&cobra.Command{
		Use:          "get",
		Short:        "Get group volume",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			v, err := c.GetGroupVolume(cmd.Context())
			if err != nil {
				return err
			}
			if flags.JSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(map[string]int{"volume": v})
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), v)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:          "set <0-100>",
		Short:        "Set group volume",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			v, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			return c.SetGroupVolume(cmd.Context(), v)
		},
	})

	return cmd
}

func newGroupMuteCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mute",
		Short: "Get or set group mute",
		Long:  "Controls GroupRenderingControl group mute on the group coordinator.",
	}

	cmd.AddCommand(&cobra.Command{
		Use:          "get",
		Short:        "Get group mute",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			m, err := c.GetGroupMute(cmd.Context())
			if err != nil {
				return err
			}
			if flags.JSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(map[string]bool{"mute": m})
			}
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), m)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:          "on",
		Short:        "Mute the whole group",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			return c.SetGroupMute(cmd.Context(), true)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:          "off",
		Short:        "Unmute the whole group",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			return c.SetGroupMute(cmd.Context(), false)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:          "toggle",
		Short:        "Toggle group mute",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			m, err := c.GetGroupMute(cmd.Context())
			if err != nil {
				return err
			}
			return c.SetGroupMute(cmd.Context(), !m)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:          "set <on|off|true|false|1|0>",
		Short:        "Set group mute",
		Args:         cobra.ExactArgs(1),
		Hidden:       true,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateTarget(flags); err != nil {
				return err
			}
			val := args[0]
			var mute bool
			switch val {
			case "on", "true", "1":
				mute = true
			case "off", "false", "0":
				mute = false
			default:
				return errors.New("invalid value: " + val)
			}
			c, err := newGroupAudioClient(cmd.Context(), flags)
			if err != nil {
				return err
			}
			return c.SetGroupMute(cmd.Context(), mute)
		},
	})

	return cmd
}
