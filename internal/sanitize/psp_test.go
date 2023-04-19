package sanitize

import (
	"testing"
	"context"

	"github.com/danielpickens/hercules/internal/cache"
	"github.com/danielpickens/hercules/internal/issues"
	"github.com/stretchr/testify/assert"
	polv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPSPSanitize(t *testing.T) {
	uu := map[string]struct {
		lister PodSecurityPolicyLister
		issues issues.Issues
	}{
		"good": {
			lister: makePSPLister("psp", pspOpts{
				rev: "policy/v1beta1",
			}),
			issues: issues.Issues{},
		},
		"deprecated": {
			lister: makePSPLister("psp", pspOpts{
				rev: "extensions/v1beta1",
			}),
			issues: issues.Issues{
				issues.Issue{
					GVR:     "policy/v1beta1/podsecuritypolicies",
					Group:   "__root__",
					Level:   2,
					Message: `[POP-403] Deprecated PodSecurityPolicy API group "extensions/v1beta1". Use "policy/v1beta1" instead`},
			},
		},
	}

	ctx := makeContext("policy/v1beta1/podsecuritypolicies", "psp")
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			psp := NewPodSecurityPolicy(issues.NewCollector(loadCodes(t), makeConfig(t)), u.lister)

			assert.Nil(t, psp.Sanitize(ctx))
			assert.Equal(t, u.issues, psp.Outcome()["default/psp"])
		})
	}
}

func TestPodCheckForMultiplePdbMatches(t *testing.T) {
	type fields struct {
		Collector   *issues.Collector
		PodMXLister PodMXLister
	}
	type args struct {
		ctx          context.Context
		podLabels    map[string]string
		podNamespace string
		pdbs         map[string]*polv1beta1.PodDisruptionBudget
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   issues.Issues
	}{
		{
			name: "pod with one label - two pdb matches",
			args: args{
				podNamespace: "namespace-1",
				podLabels:    map[string]string{"app": "test"},
				pdbs: map[string]*polv1beta1.PodDisruptionBudget{
					"pdb": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-1",
							Namespace: "namespace-1",
						},
					},
					"pdb2": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-2",
							Namespace: "namespace-1",
						},
					},
				},
			},
			want: issues.Issues{
				issues.Issue{
					GVR:     "v1/pods",
					Group:   "__root__",
					Level:   2,
					Message: "[POP-209] Pod is managed by multiple PodDisruptionBudgets (pdb-1, pdb-2)"},
			},
		},
		{
			name: "pod with one label - three pdbs - only two in pod namespace - expecting two matches",
			args: args{
				podNamespace: "namespace-1",
				podLabels:    map[string]string{"app": "test"},
				pdbs: map[string]*polv1beta1.PodDisruptionBudget{
					"pdb": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-1",
							Namespace: "namespace-1",
						},
					},
					"pdb2": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-2",
							Namespace: "namespace-1",
						},
					},
					"pdb3": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-3",
							Namespace: "namespace-2",
						},
					},
				},
			},
			want: issues.Issues{
				issues.Issue{
					GVR:     "v1/pods",
					Group:   "__root__",
					Level:   2,
					Message: "[POP-209] Pod is managed by multiple PodDisruptionBudgets (pdb-1, pdb-2)"},
			},
		},
		{
			name: "one pdb match, no issue expected",
			args: args{
				podNamespace: "namespace-1",
				podLabels:    map[string]string{"app": "test", "app2": "test2"},
				pdbs: map[string]*polv1beta1.PodDisruptionBudget{
					"pdb": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test", "app2": "test2"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-1",
							Namespace: "namespace-1",
						},
					},
					"pdb2": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app3": "test3"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "pdb-2",
							Namespace: "namespace-1",
						},
					},
				},
			},
			want: issues.Issues(nil),
		},
		{
			name: "pod with no label - no issue expected",
			args: args{
				podLabels: map[string]string{},
				pdbs: map[string]*polv1beta1.PodDisruptionBudget{
					"pdb": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "pdb-1"},
					},
					"pdb2": {
						Spec: polv1beta1.PodDisruptionBudgetSpec{
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{"app": "test"},
							},
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "pdb-2"},
					},
				},
			},
			want: issues.Issues(nil),
		},
	}
	ctx := makeContext("v1/pods", "po")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPod(issues.NewCollector(loadCodes(t), makeConfig(t)), tt.fields.PodMXLister)

			p.checkForMultiplePdbMatches(ctx, tt.args.podNamespace, tt.args.podLabels, tt.args.pdbs)
			assert.Equal(t, tt.want, p.Outcome()[""])
		})
	}
}

// Helpers...

type (
	pspOpts struct {
		rev string
	}

	psp struct {
		name string
		opts pspOpts
	}
)

func makePSPLister(n string, opts pspOpts) *psp {
	return &psp{
		name: n,
		opts: opts,
	}
}

func (r *psp) ListPodSecurityPolicies() map[string]*polv1beta1.PodSecurityPolicy {
	return map[string]*polv1beta1.PodSecurityPolicy{
		cache.FQN("default", r.name): makePSP(r.name, r.opts),
	}
}

func makePSP(n string, o pspOpts) *polv1beta1.PodSecurityPolicy {
	return &polv1beta1.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n,
			Namespace: "default",
			SelfLink:  "/api/" + o.rev,
		},
		Spec: polv1beta1.PodSecurityPolicySpec{},
	}
}
