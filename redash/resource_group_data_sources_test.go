package redash

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

type groupDataSourceBlock struct {
	id       *int // nil: unknown at plan time, the data source is created in the same apply
	viewOnly bool
}

// planGroupDataSources plans the creation of a redash_group_data_sources
// with the given data_source blocks, the way the SDK's PlanResourceChange
// does: the legacy config is shimmed from the proposed new state, where the
// computed name is unknown, and the prior state is null.
func planGroupDataSources(t *testing.T, blocks []groupDataSourceBlock) error {
	t.Helper()

	configBlocks := make([]cty.Value, 0, len(blocks))
	proposedBlocks := make([]cty.Value, 0, len(blocks))
	for _, b := range blocks {
		id := cty.UnknownVal(cty.Number)
		if b.id != nil {
			id = cty.NumberIntVal(int64(*b.id))
		}
		configBlocks = append(configBlocks, cty.ObjectVal(map[string]cty.Value{
			"data_source_id": id,
			"name":           cty.NullVal(cty.String),
			"view_only":      cty.BoolVal(b.viewOnly),
		}))
		proposedBlocks = append(proposedBlocks, cty.ObjectVal(map[string]cty.Value{
			"data_source_id": id,
			"name":           cty.UnknownVal(cty.String),
			"view_only":      cty.BoolVal(b.viewOnly),
		}))
	}

	config := cty.ObjectVal(map[string]cty.Value{
		"id":          cty.NullVal(cty.String),
		"group_id":    cty.NumberIntVal(21),
		"data_source": cty.SetVal(configBlocks),
	})
	proposed := cty.ObjectVal(map[string]cty.Value{
		"id":          cty.UnknownVal(cty.String),
		"group_id":    cty.NumberIntVal(21),
		"data_source": cty.SetVal(proposedBlocks),
	})

	r := resourceGroupDataSources()
	ty := r.CoreConfigSchema().ImpliedType()
	prior, err := r.ShimInstanceStateFromValue(cty.NullVal(ty))
	require.NoError(t, err)
	prior.RawState = cty.NullVal(ty)
	prior.RawPlan = proposed
	prior.RawConfig = config

	_, err = r.SimpleDiff(context.Background(), prior, terraform.NewResourceConfigShimmed(proposed, r.CoreConfigSchema()), nil)

	return err
}

func TestResourceGroupDataSources_planRejectsDuplicates(t *testing.T) {
	known := func(id int) *int { return &id }

	tests := []struct {
		name    string
		blocks  []groupDataSourceBlock
		wantErr string
	}{
		{
			name:   "ids unknown at plan time",
			blocks: []groupDataSourceBlock{{id: nil, viewOnly: false}, {id: nil, viewOnly: true}},
		},
		{
			name:   "unknown id next to a known one",
			blocks: []groupDataSourceBlock{{id: nil, viewOnly: true}, {id: known(4), viewOnly: false}},
		},
		{
			name:   "distinct known ids",
			blocks: []groupDataSourceBlock{{id: known(4)}, {id: known(5)}},
		},
		{
			name:    "duplicate known ids",
			blocks:  []groupDataSourceBlock{{id: known(4), viewOnly: false}, {id: known(4), viewOnly: true}},
			wantErr: "duplicate data_source_id 4",
		},
		{
			name:    "duplicate known ids next to an unknown one",
			blocks:  []groupDataSourceBlock{{id: nil}, {id: known(4), viewOnly: false}, {id: known(4), viewOnly: true}},
			wantErr: "duplicate data_source_id 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := planGroupDataSources(t, tt.blocks)

			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

func TestResourceGroupDataSources_planWithoutRawPlan(t *testing.T) {
	proposed := cty.ObjectVal(map[string]cty.Value{
		"id":       cty.UnknownVal(cty.String),
		"group_id": cty.NumberIntVal(21),
		"data_source": cty.SetVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"data_source_id": cty.NumberIntVal(4),
			"name":           cty.UnknownVal(cty.String),
			"view_only":      cty.False,
		})}),
	})

	r := resourceGroupDataSources()
	ty := r.CoreConfigSchema().ImpliedType()
	prior, err := r.ShimInstanceStateFromValue(cty.NullVal(ty))
	require.NoError(t, err)
	prior.RawState = cty.NullVal(ty)
	prior.RawPlan = cty.NullVal(ty)
	prior.RawConfig = cty.NullVal(ty)

	require.NotPanics(t, func() {
		_, err = r.SimpleDiff(context.Background(), prior, terraform.NewResourceConfigShimmed(proposed, r.CoreConfigSchema()), nil)
	})
	require.NoError(t, err)
}
