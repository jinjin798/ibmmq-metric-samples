package ibmmq

/*
  Copyright (c) IBM Corporation 2016

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

   Contributors:
     Mark Taylor - Initial Contribution
*/

/*

#include <stdlib.h>
#include <string.h>
#include <cmqc.h>
#include <cmqxc.h>

void setScoKeyRepoPassword(MQSCO *mqsco, PMQCHAR keyRepoPassword, MQLONG length) {
#if defined(MQSCO_VERSION_6) && MQSCO_CURRENT_VERSION >= MQSCO_VERSION_6
  if (mqsco->Version < MQSCO_VERSION_6) {
	  mqsco->Version = MQSCO_VERSION_6;
  }
  mqsco->KeyRepoPasswordOffset = 0;
  mqsco->KeyRepoPasswordPtr = NULL;
  mqsco->KeyRepoPasswordLength = length;
  if (keyRepoPassword != NULL && length > 0) {
    mqsco->KeyRepoPasswordPtr = keyRepoPassword;
  }
#else
  if (keyRepoPassword != NULL) {
    free(keyRepoPassword);
  }
#endif
}

void freeScoKeyRepoPassword(MQSCO *mqsco) {
#if defined(MQSCO_VERSION_6) && MQSCO_CURRENT_VERSION >= MQSCO_VERSION_6
  if (mqsco->KeyRepoPasswordPtr != NULL) {
    free(mqsco->KeyRepoPasswordPtr);
  }
#endif
}

*/
import "C"

/*
MQSCO is a structure containing the MQ SSL/TLS Configuration (MQSCO)
options.
*/
type MQSCO struct {
	KeyRepository          string
	CryptoHardware         string
	KeyResetCount          int32
	FipsRequired           bool
	EncryptionPolicySuiteB [4]int32
	CertificateValPolicy   int32
	CertificateLabel       string
	KeyRepoPassword        string
}

/*
NewMQSCO fills in default values for the MQSCO structure
*/
func NewMQSCO() *MQSCO {

	sco := new(MQSCO)

	sco.KeyRepository = ""
	sco.CryptoHardware = ""
	sco.KeyResetCount = int32(C.MQSCO_RESET_COUNT_DEFAULT)
	sco.FipsRequired = false
	sco.EncryptionPolicySuiteB[0] = int32(C.MQ_SUITE_B_NONE)
	for i := 1; i < 4; i++ {
		sco.EncryptionPolicySuiteB[i] = int32(C.MQ_SUITE_B_NOT_AVAILABLE)
	}
	sco.CertificateValPolicy = int32(C.MQ_CERT_VAL_POLICY_DEFAULT)
	sco.CertificateLabel = ""
	sco.KeyRepoPassword = ""

	return sco
}

/*
It is expected that copyXXtoC and copyXXfromC will be called as
matching pairs.
*/
func copySCOtoC(mqsco *C.MQSCO, gosco *MQSCO) {

	setMQIString((*C.char)(&mqsco.StrucId[0]), "SCO ", 4)
	mqsco.Version = C.MQSCO_VERSION_5
	setMQIString((*C.char)(&mqsco.KeyRepository[0]), gosco.KeyRepository, C.MQ_SSL_KEY_REPOSITORY_LENGTH)
	setMQIString((*C.char)(&mqsco.CryptoHardware[0]), gosco.CryptoHardware, C.MQ_SSL_CRYPTO_HARDWARE_LENGTH)
	mqsco.AuthInfoRecCount = 0
	mqsco.AuthInfoRecOffset = 0
	mqsco.AuthInfoRecPtr = nil
	mqsco.KeyResetCount = C.MQLONG(gosco.KeyResetCount)
	if gosco.FipsRequired {
		mqsco.FipsRequired = C.MQSSL_FIPS_YES
	} else {
		mqsco.FipsRequired = C.MQSSL_FIPS_NO
	}
	for i := 0; i < 4; i++ {
		mqsco.EncryptionPolicySuiteB[i] = C.MQLONG(gosco.EncryptionPolicySuiteB[i])
	}
	mqsco.CertificateValPolicy = C.MQLONG(gosco.CertificateValPolicy)
	setMQIString((*C.char)(&mqsco.CertificateLabel[0]), gosco.CertificateLabel, C.MQ_CERT_LABEL_LENGTH)

	if gosco.KeyRepoPassword != "" {
		C.setScoKeyRepoPassword(mqsco, C.PMQCHAR(C.CString(gosco.KeyRepoPassword)), C.MQLONG(len(gosco.KeyRepoPassword)))
	}
	return
}

/*
Only need to free the KeyRepoPassword if it were set
*/
func copySCOfromC(mqsco *C.MQSCO, gosco *MQSCO) {
	C.freeScoKeyRepoPassword(mqsco) // The code in this function checks validity

	return
}
