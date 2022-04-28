package httpx

import "net/http"

/*
创建人员：云深不知处
创建时间：2022/1/2
程序功能：
*/

// TLSData contains the relevant Transport Layer Security information
type TLSData struct {
	DNSNames           []string `json:"dns_names,omitempty"`
	EmailAddresses     []string `json:"email_addresses,omitempty"`
	CommonName         []string `json:"common_name,omitempty"`
	Organization       []string `json:"organization,omitempty"`
	IssuerCommonName   []string `json:"issuer_common_name,omitempty"`
	IssuerOrg          []string `json:"issuer_organization,omitempty"`
	OrganizationalUnit []string `json:"organizational_unit,omitempty"`
	Issuer             []string `json:"issuer,omitempty"`
	Subject            []string `json:"subject,omitempty"`
}

// TLSGrab fills the TLSData
func (h *HTTPX) TLSGrab(r *http.Response) *TLSData {
	if r.TLS != nil {
		var tlsdata TLSData
		for _, certificate := range r.TLS.PeerCertificates {
			tlsdata.OrganizationalUnit = append(tlsdata.OrganizationalUnit, certificate.Subject.OrganizationalUnit...)
			tlsdata.DNSNames = append(tlsdata.DNSNames, certificate.DNSNames...)
			tlsdata.EmailAddresses = append(tlsdata.EmailAddresses, certificate.EmailAddresses...)
			tlsdata.CommonName = append(tlsdata.CommonName, certificate.Subject.CommonName)
			tlsdata.Organization = append(tlsdata.Organization, certificate.Subject.Organization...)
			tlsdata.IssuerOrg = append(tlsdata.IssuerOrg, certificate.Issuer.Organization...)
			tlsdata.IssuerCommonName = append(tlsdata.IssuerCommonName, certificate.Issuer.CommonName)
			tlsdata.Subject = append(tlsdata.Subject, certificate.Subject.String())
			tlsdata.Issuer = append(tlsdata.Issuer, certificate.Issuer.String())
		}
		return &tlsdata
	}
	return nil
}
