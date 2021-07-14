package provider

import (
	"context"
	"strings"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel"
	"github.com/chronark/terraform-provider-vercel/pkg/vercel/project_domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectDomain() *schema.Resource {
	return &schema.Resource{
		Description: "https://vercel.com/docs/api#endpoints/projects",

		CreateContext: resourceProjectDomainCreate,
		ReadContext:   resourceProjectDomainRead,
		UpdateContext: resourceProjectDomainUpdate,
		DeleteContext: resourceProjectDomainDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "Internal id of project to manage the domains of",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the domain to add.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"team_id": {
				Description: "By default, you can access resources contained within your own user account. To access resources owned by a team, you can pass in the team ID",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			"redirect": {
				Description: "Target destination domain for redirect.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"redirect_status_code": {
				Description: "The redirect status code (301, 302, 307, 308).",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"git_branch": {
				Description: "Git branch for the domain to be auto assigned to. The Project's production branch is the default.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"created_at": {
				Description: "A number containing the date when the project domain was created in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"updated_at": {
				Description: "A number containing the date when the project domain was updated in milliseconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceProjectDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	projectDomain := project_domain.CreateProjectDomain{
		Name: d.Get("name").(string),
	}

	if redirect, ok := d.GetOk("redirect"); ok {
		projectDomain.Redirect = redirect.(string)
	}

	if redirectStatusCode, ok := d.GetOk("redirect_status_code"); ok {
		projectDomain.RedirectStatusCode = redirectStatusCode.(int)
	}

	if gitBranch, ok := d.GetOk("git_branch"); ok {
		projectDomain.GitBranch = gitBranch.(string)
	}

	id, err := client.ProjectDomain.Create(d.Get("project_id").(string), projectDomain, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("redirect", projectDomain.Redirect)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("redirect_status_code", projectDomain.RedirectStatusCode)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceProjectDomainRead(ctx, d, meta)
}

func resourceProjectDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)

	id := d.Id()

	subs := strings.SplitN(id, ":", 2)
	projectID, name := subs[0], subs[1]

	projectDomain, err := client.ProjectDomain.Read(projectID, name, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", projectDomain.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("project_id", projectDomain.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", projectDomain.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated_at", projectDomain.UpdatedAt)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceProjectDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)
	var update project_domain.UpdateProjectDomain

	id := d.Id()
	subs := strings.SplitN(id, ":", 2)
	projectID, name := subs[0], subs[1]

	if d.HasChange("redirect") {
		update.Redirect = d.Get("redirect").(string)
	}

	projectDomain, err := client.ProjectDomain.Update(projectID, name, update, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", projectDomain.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("project_id", projectDomain.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("created_at", projectDomain.CreatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("updated_at", projectDomain.UpdatedAt)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("redirect", projectDomain.Redirect)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("redirect_status_code", 307)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func resourceProjectDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*vercel.Client)
	id := d.Id()
	subs := strings.SplitN(id, ":", 2)
	projectID, name := subs[0], subs[1]

	err := client.ProjectDomain.Delete(projectID, name, d.Get("team_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}
