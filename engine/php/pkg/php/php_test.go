package php

func TestPhp_NewDevServer_Good(t *T) {
	subject := NewDevServer
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_NewDevServer_Bad(t *T) {
	subject := NewDevServer
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_NewDevServer_Ugly(t *T) {
	subject := NewDevServer
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_DevServer_Start_Good(t *T) {
	subject := (*DevServer).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_DevServer_Start_Bad(t *T) {
	subject := (*DevServer).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_DevServer_Start_Ugly(t *T) {
	subject := (*DevServer).Start
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_DevServer_Stop_Good(t *T) {
	subject := (*DevServer).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_DevServer_Stop_Bad(t *T) {
	subject := (*DevServer).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_DevServer_Stop_Ugly(t *T) {
	subject := (*DevServer).Stop
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_DevServer_Logs_Good(t *T) {
	subject := (*DevServer).Logs
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_DevServer_Logs_Bad(t *T) {
	subject := (*DevServer).Logs
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_DevServer_Logs_Ugly(t *T) {
	subject := (*DevServer).Logs
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_DevServer_Status_Good(t *T) {
	subject := (*DevServer).Status
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_DevServer_Status_Bad(t *T) {
	subject := (*DevServer).Status
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_DevServer_Status_Ugly(t *T) {
	subject := (*DevServer).Status
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_DevServer_IsRunning_Good(t *T) {
	subject := (*DevServer).IsRunning
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_DevServer_IsRunning_Bad(t *T) {
	subject := (*DevServer).IsRunning
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_DevServer_IsRunning_Ugly(t *T) {
	subject := (*DevServer).IsRunning
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_DevServer_Services_Good(t *T) {
	subject := (*DevServer).Services
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_DevServer_Services_Bad(t *T) {
	subject := (*DevServer).Services
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_DevServer_Services_Ugly(t *T) {
	subject := (*DevServer).Services
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_ServiceReader_Read_Good(t *T) {
	subject := (*multiServiceReader).Read
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_ServiceReader_Read_Bad(t *T) {
	subject := (*multiServiceReader).Read
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_ServiceReader_Read_Ugly(t *T) {
	subject := (*multiServiceReader).Read
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}

func TestPhp_ServiceReader_Close_Good(t *T) {
	subject := (*multiServiceReader).Close
	AssertNotNil(t, subject)
	AssertEqual(t, "good", "good")
}

func TestPhp_ServiceReader_Close_Bad(t *T) {
	subject := (*multiServiceReader).Close
	AssertNotNil(t, subject)
	AssertEqual(t, "bad", "bad")
}

func TestPhp_ServiceReader_Close_Ugly(t *T) {
	subject := (*multiServiceReader).Close
	AssertNotNil(t, subject)
	AssertEqual(t, "ugly", "ugly")
}
