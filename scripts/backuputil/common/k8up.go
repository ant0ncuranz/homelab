package common

import (
	"context"
	"encoding/json"
	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	k8upv1a1 "github.com/vshn/k8up/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func CreateRestore(client *kubernetes.Clientset, namespace string, restoreName string, pvcName string, snapshot string) error {
	restore := k8upv1a1.Restore{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Restore",
			APIVersion: "k8up.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      restoreName,
			Namespace: namespace,
		},
		Spec: k8upv1a1.RestoreSpec{
			RestoreMethod: &k8upv1a1.RestoreMethod{
				Folder: &k8upv1a1.FolderRestore{
					&v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvcName,
					},
				},
			},
			Snapshot: snapshot,
		},
	}

	body, err := json.Marshal(restore)
	if err != nil {
		return err
	}

	rsp := k8upv1a1.Restore{}
	absPath := "/apis/k8up.io/v1/namespaces/" + namespace + "/restores"
	return client.RESTClient().Post().AbsPath(absPath).Body(body).Do(context.Background()).Into(&rsp)
}

func CreateBackup(client *kubernetes.Clientset, namespace string, backupName string) error {
	backup := k8upv1a1.Backup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Backup",
			APIVersion: "k8up.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupName,
			Namespace: namespace,
		},
	}

	body, err := json.Marshal(backup)
	if err != nil {
		return err
	}

	rsp := k8upv1a1.Backup{}
	absPath := "/apis/k8up.io/v1/namespaces/" + namespace + "/backups"
	return client.RESTClient().Post().AbsPath(absPath).Body(body).Do(context.Background()).Into(&rsp)
}

func CreateDatabaseBackup(client *kubernetes.Clientset, backupName string) error {
	backup := cnpgv1.Backup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Backup",
			APIVersion: "postgresql.cnpg.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      backupName,
			Namespace: "database",
		},
		Spec: cnpgv1.BackupSpec{
			Cluster: cnpgv1.LocalObjectReference{
				Name: "database-cluster",
			},
		},
	}

	body, err := json.Marshal(backup)
	if err != nil {
		return err
	}

	rsp := cnpgv1.Backup{}
	absPath := "/apis/postgresql.cnpg.io/v1/namespaces/database/backups"
	return client.RESTClient().Post().AbsPath(absPath).Body(body).Do(context.Background()).Into(&rsp)
}

func RestorePvc(client *kubernetes.Clientset, namespace string, snapshot string, pvcName string, pvcStorage string, storageClass string) error {
	if err := CreatePvc(client, namespace, pvcName, pvcStorage, storageClass); err != nil {
		return err
	}

	restoreName := "restore-" + pvcName
	if err := CreateRestore(client, namespace, restoreName, pvcName, snapshot); err != nil {
		log.Fatal(err)
	}

	return nil
}
