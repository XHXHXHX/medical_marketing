package server

import "github.com/XHXHXHX/medical_marketing/service/user"

var RoleInterfaceMap = map[string]map[user.Role]struct{}{
	"/v1api.report.apiReportService/UserCreate": {
		user.RoleMarketManager:{},
		user.RoleCustomManager: {},
	},
	"/v1api.report.apiReportService/UserChangeStatus": {
		user.RoleMarketManager:{},
		user.RoleCustomManager: {},
	},


	"/v1api.report.apiReportService/ReportList": {
		user.RoleMarketManager:{},
		user.RoleMarketStaff: {},
	},
	"/v1api.report.apiReportService/ReportCreate": {
		user.RoleMarketManager:{},
		user.RoleMarketStaff: {},
	},
	"/v1api.report.apiReportService/ReportRecover": {
		user.RoleMarketManager:{},
		user.RoleMarketStaff: {},
	},


	"/v1api.customer.apiCustomerServerService/CustomerServerDistribute": {
		user.RoleCustomManager: {},
	},
	"/v1api.customer.apiCustomerServerService/CustomerServerList": {
		user.RoleCustomManager: {},
		user.RoleCustomStaff: {},
	},
	"/v1api.customer.apiCustomerServerService/CustomerServerResult": {
		user.RoleCustomManager: {},
		user.RoleCustomStaff: {},
	},


	"/v1api.statistics.apiStatisticsService/StatisticsMarket": {
		user.RoleMarketManager: {},
	},
	"/v1api.statistics.apiStatisticsService/StatisticsCustomer": {
		user.RoleCustomManager: {},
	},
}

var SkipInterfaceMap = map[string]struct{}{
	"/v1api.report.apiUserService/Login": {},
}