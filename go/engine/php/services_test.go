package php

func TestServices_Service_Name_Good(t *T) {
	subject := (*baseService).Name
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_Service_Name_Bad(t *T) {
	subject := (*baseService).Name
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_Service_Name_Ugly(t *T) {
	subject := (*baseService).Name
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_Service_Status_Good(t *T) {
	subject := (*baseService).Status
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_Service_Status_Bad(t *T) {
	subject := (*baseService).Status
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_Service_Status_Ugly(t *T) {
	subject := (*baseService).Status
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_Service_Logs_Good(t *T) {
	subject := (*baseService).Logs
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_Service_Logs_Bad(t *T) {
	subject := (*baseService).Logs
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_Service_Logs_Ugly(t *T) {
	subject := (*baseService).Logs
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_NewFrankenPHPService_Good(t *T) {
	subject := NewFrankenPHPService
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_NewFrankenPHPService_Bad(t *T) {
	subject := NewFrankenPHPService
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_NewFrankenPHPService_Ugly(t *T) {
	subject := NewFrankenPHPService
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_FrankenPHPService_Start_Good(t *T) {
	subject := (*FrankenPHPService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_FrankenPHPService_Start_Bad(t *T) {
	subject := (*FrankenPHPService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_FrankenPHPService_Start_Ugly(t *T) {
	subject := (*FrankenPHPService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_FrankenPHPService_Stop_Good(t *T) {
	subject := (*FrankenPHPService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_FrankenPHPService_Stop_Bad(t *T) {
	subject := (*FrankenPHPService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_FrankenPHPService_Stop_Ugly(t *T) {
	subject := (*FrankenPHPService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_NewViteService_Good(t *T) {
	subject := NewViteService
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_NewViteService_Bad(t *T) {
	subject := NewViteService
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_NewViteService_Ugly(t *T) {
	subject := NewViteService
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_ViteService_Start_Good(t *T) {
	subject := (*ViteService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_ViteService_Start_Bad(t *T) {
	subject := (*ViteService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_ViteService_Start_Ugly(t *T) {
	subject := (*ViteService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_ViteService_Stop_Good(t *T) {
	subject := (*ViteService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_ViteService_Stop_Bad(t *T) {
	subject := (*ViteService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_ViteService_Stop_Ugly(t *T) {
	subject := (*ViteService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_NewHorizonService_Good(t *T) {
	subject := NewHorizonService
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_NewHorizonService_Bad(t *T) {
	subject := NewHorizonService
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_NewHorizonService_Ugly(t *T) {
	subject := NewHorizonService
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_HorizonService_Start_Good(t *T) {
	subject := (*HorizonService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_HorizonService_Start_Bad(t *T) {
	subject := (*HorizonService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_HorizonService_Start_Ugly(t *T) {
	subject := (*HorizonService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_HorizonService_Stop_Good(t *T) {
	subject := (*HorizonService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_HorizonService_Stop_Bad(t *T) {
	subject := (*HorizonService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_HorizonService_Stop_Ugly(t *T) {
	subject := (*HorizonService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_NewReverbService_Good(t *T) {
	subject := NewReverbService
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_NewReverbService_Bad(t *T) {
	subject := NewReverbService
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_NewReverbService_Ugly(t *T) {
	subject := NewReverbService
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_ReverbService_Start_Good(t *T) {
	subject := (*ReverbService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_ReverbService_Start_Bad(t *T) {
	subject := (*ReverbService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_ReverbService_Start_Ugly(t *T) {
	subject := (*ReverbService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_ReverbService_Stop_Good(t *T) {
	subject := (*ReverbService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_ReverbService_Stop_Bad(t *T) {
	subject := (*ReverbService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_ReverbService_Stop_Ugly(t *T) {
	subject := (*ReverbService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_NewRedisService_Good(t *T) {
	subject := NewRedisService
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_NewRedisService_Bad(t *T) {
	subject := NewRedisService
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_NewRedisService_Ugly(t *T) {
	subject := NewRedisService
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_RedisService_Start_Good(t *T) {
	subject := (*RedisService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_RedisService_Start_Bad(t *T) {
	subject := (*RedisService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_RedisService_Start_Ugly(t *T) {
	subject := (*RedisService).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_RedisService_Stop_Good(t *T) {
	subject := (*RedisService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_RedisService_Stop_Bad(t *T) {
	subject := (*RedisService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_RedisService_Stop_Ugly(t *T) {
	subject := (*RedisService).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_Reader_Read_Good(t *T) {
	subject := (*tailReader).Read
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_Reader_Read_Bad(t *T) {
	subject := (*tailReader).Read
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_Reader_Read_Ugly(t *T) {
	subject := (*tailReader).Read
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestServices_Reader_Close_Good(t *T) {
	subject := (*tailReader).Close
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestServices_Reader_Close_Bad(t *T) {
	subject := (*tailReader).Close
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestServices_Reader_Close_Ugly(t *T) {
	subject := (*tailReader).Close
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
