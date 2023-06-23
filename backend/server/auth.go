package server

import "github.com/XHXHXHX/medical_marketing/service/user"

var RoleInterfaceMap = map[string]map[user.Role]struct{}{
	"/v1api.report.apiReportService/UserCreate": {
		user.RoleMarketManager:{},
		user.RoleCustomManager: {},
		user.RoleFinance: {},
	},
	"/v1api.report.apiReportService/UserChangeStatus": {
		user.RoleMarketManager:{},
		user.RoleCustomManager: {},
		user.RoleFinance: {},
	},


	"/v1api.report.apiReportService/ReportList": {
		user.RoleMarketManager:{},
		user.RoleMarketStaff: {},
		user.RoleCustomManager: {},
		user.RoleCustomStaff: {},
		user.RoleFinance: {},
	},
	"/v1api.report.apiReportService/ReportCreate": {
		user.RoleMarketManager:{},
		user.RoleMarketStaff: {},
		user.RoleCustomManager: {},
		user.RoleCustomStaff: {},
		user.RoleFinance: {},
	},
	"/v1api.report.apiReportService/ReportRecover": {
		user.RoleMarketManager:{},
		user.RoleMarketStaff: {},
		user.RoleFinance: {},
	},
	"/v1api.report.apiReportService/ReportChangeMatch": {
		user.RoleFinance:{},
	},


	"/v1api.customer.apiCustomerServerService/CustomerServerDistribute": {
		user.RoleCustomManager: {},
		user.RoleFinance: {},
	},
	"/v1api.customer.apiCustomerServerService/CustomerServerList": {
		user.RoleCustomManager: {},
		user.RoleCustomStaff: {},
		user.RoleFinance: {},
	},
	"/v1api.customer.apiCustomerServerService/CustomerServerResult": {
		user.RoleCustomManager: {},
		user.RoleCustomStaff: {},
		user.RoleFinance: {},
	},


	"/v1api.statistics.apiStatisticsService/StatisticsMarket": {
		user.RoleMarketManager: {},
		user.RoleFinance: {},
	},
	"/v1api.statistics.apiStatisticsService/StatisticsCustomer": {
		user.RoleCustomManager: {},
		user.RoleFinance: {},
	},
}

var SkipInterfaceMap = map[string]struct{}{
	"/v1api.report.apiUserService/Login": {},
}