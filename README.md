# canal-sla

This is the branch for the SLA part of the Canal

This version is for the restore when the code crashes

When the code restart after crush, it will first load the etcd configuration. However, before configing, we delete all the existing configuration by setting all the action to Delete. And then change all the action to ADD except the rule with the "delete" action.






