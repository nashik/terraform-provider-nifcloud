package securitygrouprule

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nifcloud/nifcloud-sdk-go/nifcloud"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
)

func flatten(d *schema.ResourceData, res *computing.DescribeSecurityGroupsResponse) error {
	if res == nil || len(res.SecurityGroupInfo) == 0 {
		d.SetId("")
		return nil
	}

	names := make([]string, len(res.SecurityGroupInfo))
	for i, s := range res.SecurityGroupInfo {
		names[i] = nifcloud.StringValue(s.GroupName)
	}

	var rule *computing.IpPermissions
	r := expandAuthorizeSecurityGroupIngressInputList(d)[0].IpPermissions[0]
	for i := range res.SecurityGroupInfo[0].IpPermissions {
		p := &res.SecurityGroupInfo[0].IpPermissions[i]

		if p.ToPort != nil && r.ToPort != nil && nifcloud.Int64Value(p.ToPort) != nifcloud.Int64Value(r.ToPort) {
			continue
		}

		if p.FromPort != nil && r.FromPort != nil && nifcloud.Int64Value(p.FromPort) != nifcloud.Int64Value(r.FromPort) {
			continue
		}

		if nifcloud.StringValue(p.IpProtocol) != string(r.IpProtocol) {
			continue
		}

		if nifcloud.StringValue(p.InOut) != string(r.InOut) {
			continue
		}

		findCidrIP := false
		if len(r.ListOfRequestIpRanges) > 0 {
			for _, ip := range p.IpRanges {
				if nifcloud.StringValue(ip.CidrIp) == nifcloud.StringValue(r.ListOfRequestIpRanges[0].CidrIp) {
					findCidrIP = true
					break
				}
			}
		}
		if findCidrIP {
			rule = p
			break
		}

		findGroup := false
		if len(r.ListOfRequestGroups) > 0 {
			for _, gn := range p.Groups {
				if nifcloud.StringValue(gn.GroupName) == nifcloud.StringValue(r.ListOfRequestGroups[0].GroupName) {
					findGroup = true
					break
				}
			}
		}
		if findGroup {
			rule = p
			break
		}
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("security_group_names", names); err != nil {
		return err
	}

	if err := d.Set("type", rule.InOut); err != nil {
		return err
	}

	if len(rule.IpRanges) > 0 {
		if err := d.Set("cidr_ip", rule.IpRanges[0].CidrIp); err != nil {
			return err
		}
	}

	if rule.FromPort != nil {
		if err := d.Set("from_port", rule.FromPort); err != nil {
			return err
		}
	}

	if err := d.Set("protocol", rule.IpProtocol); err != nil {
		return err
	}

	if len(rule.Groups) > 0 {
		if err := d.Set("source_security_group_name", rule.Groups[0].GroupName); err != nil {
			return err
		}
	}

	if rule.ToPort != nil {
		if rule.FromPort == rule.ToPort && d.Get("to_port") != nil {
			if err := d.Set("to_port", rule.ToPort); err != nil {
				return err
			}
		}
	}

	if rule.Description != nil {
		if err := d.Set("description", rule.Description); err != nil {
			return err
		}
	}

	if strings.Contains(d.Id(), "_") {
		// import so fix the id
		id := idHash(expandAuthorizeSecurityGroupIngressInputList(d))
		d.SetId(id)
	}

	return nil
}
