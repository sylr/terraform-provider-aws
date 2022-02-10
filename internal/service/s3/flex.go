package s3

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func ExpandAccessControlTranslation(l []interface{}) *s3.AccessControlTranslation {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.AccessControlTranslation{}

	if v, ok := tfMap["owner"].(string); ok && v != "" {
		result.Owner = aws.String(v)
	}

	return result
}

func ExpandEncryptionConfiguration(l []interface{}) *s3.EncryptionConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.EncryptionConfiguration{}

	if v, ok := tfMap["replica_kms_key_id"].(string); ok && v != "" {
		result.ReplicaKmsKeyID = aws.String(v)
	}

	return result
}

func ExpandDeleteMarkerReplication(l []interface{}) *s3.DeleteMarkerReplication {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.DeleteMarkerReplication{}

	if v, ok := tfMap["status"].(string); ok && v != "" {
		result.Status = aws.String(v)
	}

	return result
}

func ExpandDestination(l []interface{}) *s3.Destination {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.Destination{}

	if v, ok := tfMap["access_control_translation"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.AccessControlTranslation = ExpandAccessControlTranslation(v)
	}

	if v, ok := tfMap["account"].(string); ok && v != "" {
		result.Account = aws.String(v)
	}

	if v, ok := tfMap["bucket"].(string); ok && v != "" {
		result.Bucket = aws.String(v)
	}

	if v, ok := tfMap["encryption_configuration"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.EncryptionConfiguration = ExpandEncryptionConfiguration(v)
	}

	if v, ok := tfMap["metrics"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.Metrics = ExpandMetrics(v)
	}

	if v, ok := tfMap["replication_time"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.ReplicationTime = ExpandReplicationTime(v)
	}

	if v, ok := tfMap["storage_class"].(string); ok && v != "" {
		result.StorageClass = aws.String(v)
	}

	return result
}

func ExpandExistingObjectReplication(l []interface{}) *s3.ExistingObjectReplication {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.ExistingObjectReplication{}

	if v, ok := tfMap["status"].(string); ok && v != "" {
		result.Status = aws.String(v)
	}

	return result
}

func ExpandFilter(l []interface{}) *s3.ReplicationRuleFilter {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.ReplicationRuleFilter{}

	if v, ok := tfMap["and"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.And = ExpandReplicationRuleAndOperator(v)
	}

	if v, ok := tfMap["prefix"].(string); ok && v != "" {
		result.Prefix = aws.String(v)
	}

	if v, ok := tfMap["tag"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		tags := Tags(tftags.New(v[0]).IgnoreAWS())
		if len(tags) > 0 {
			result.Tag = tags[0]
		}
	}

	return result
}

func ExpandLifecycleRuleAbortIncompleteMultipartUpload(m map[string]interface{}) *s3.AbortIncompleteMultipartUpload {
	if len(m) == 0 {
		return nil
	}

	result := &s3.AbortIncompleteMultipartUpload{}

	if v, ok := m["days_after_initiation"].(int); ok {
		result.DaysAfterInitiation = aws.Int64(int64(v))
	}

	return result
}

func ExpandLifecycleRuleExpiration(l []interface{}) (*s3.LifecycleExpiration, error) {
	if len(l) == 0 {
		return nil, nil
	}

	result := &s3.LifecycleExpiration{}

	if l[0] == nil {
		return result, nil
	}

	m := l[0].(map[string]interface{})

	if v, ok := m["date"].(string); ok && v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return nil, fmt.Errorf("error parsing S3 Bucket Lifecycle Rule Expiration date: %w", err)
		}
		result.Date = aws.Time(t)
	}

	if v, ok := m["days"].(int); ok && v > 0 {
		result.Days = aws.Int64(int64(v))
	}

	// This cannot be specified with Days or Date
	if v, ok := m["expired_object_delete_marker"].(bool); ok && result.Date == nil && result.Days == nil {
		result.ExpiredObjectDeleteMarker = aws.Bool(v)
	}

	return result, nil
}

// ExpandLifecycleRuleFilter ensures a Filter can have only 1 of prefix, tag, or and
func ExpandLifecycleRuleFilter(l []interface{}) *s3.LifecycleRuleFilter {
	if len(l) == 0 {
		return nil
	}

	result := &s3.LifecycleRuleFilter{}

	if l[0] == nil {
		return result
	}

	m := l[0].(map[string]interface{})

	if v, ok := m["and"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.And = ExpandLifecycleRuleFilterAndOperator(v[0].(map[string]interface{}))
	}

	if v, ok := m["object_size_greater_than"].(int); ok && v > 0 {
		result.ObjectSizeGreaterThan = aws.Int64(int64(v))
	}

	if v, ok := m["object_size_less_than"].(int); ok && v > 0 {
		result.ObjectSizeLessThan = aws.Int64(int64(v))
	}

	if v, ok := m["tag"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		tags := Tags(tftags.New(v[0]).IgnoreAWS())
		if len(tags) > 0 {
			result.Tag = tags[0]
		}
	}

	if v, ok := m["prefix"].(string); ok && result.And == nil && result.Tag == nil {
		result.Prefix = aws.String(v)
	}

	return result
}

func ExpandLifecycleRuleFilterAndOperator(m map[string]interface{}) *s3.LifecycleRuleAndOperator {
	if len(m) == 0 {
		return nil
	}

	result := &s3.LifecycleRuleAndOperator{}

	if v, ok := m["object_size_greater_than"].(int); ok && v > 0 {
		result.ObjectSizeGreaterThan = aws.Int64(int64(v))
	}

	if v, ok := m["object_size_less_than"].(int); ok && v > 0 {
		result.ObjectSizeLessThan = aws.Int64(int64(v))
	}

	if v, ok := m["prefix"].(string); ok {
		result.Prefix = aws.String(v)
	}

	if v, ok := m["tags"].(map[string]interface{}); ok && len(v) > 0 {
		tags := Tags(tftags.New(v).IgnoreAWS())
		if len(tags) > 0 {
			result.Tags = tags
		}
	}

	return result
}

func ExpandLifecycleRuleNoncurrentVersionExpiration(m map[string]interface{}) *s3.NoncurrentVersionExpiration {
	if len(m) == 0 {
		return nil
	}

	result := &s3.NoncurrentVersionExpiration{}

	if v, ok := m["newer_noncurrent_versions"].(int); ok && v > 0 {
		result.NewerNoncurrentVersions = aws.Int64(int64(v))
	}

	if v, ok := m["noncurrent_days"].(int); ok {
		result.NoncurrentDays = aws.Int64(int64(v))
	}

	return result
}

func ExpandLifecycleRuleNoncurrentVersionTransitions(l []interface{}) []*s3.NoncurrentVersionTransition {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	var results []*s3.NoncurrentVersionTransition

	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		transition := &s3.NoncurrentVersionTransition{}

		if v, ok := tfMap["newer_noncurrent_versions"].(int); ok && v > 0 {
			transition.NewerNoncurrentVersions = aws.Int64(int64(v))
		}

		if v, ok := tfMap["noncurrent_days"].(int); ok {
			transition.NoncurrentDays = aws.Int64(int64(v))
		}

		if v, ok := tfMap["storage_class"].(string); ok && v != "" {
			transition.StorageClass = aws.String(v)
		}

		results = append(results, transition)
	}

	return results
}

func ExpandLifecycleRuleTransitions(l []interface{}) ([]*s3.Transition, error) {
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}

	var results []*s3.Transition

	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		transition := &s3.Transition{}

		if v, ok := tfMap["date"].(string); ok && v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("error parsing S3 Bucket Lifecycle Rule Transition date: %w", err)
			}
			transition.Date = aws.Time(t)
		}

		if v, ok := tfMap["days"].(int); ok && v > 0 {
			transition.Days = aws.Int64(int64(v))
		}

		if v, ok := tfMap["storage_class"].(string); ok && v != "" {
			transition.StorageClass = aws.String(v)
		}

		results = append(results, transition)
	}

	return results, nil
}

func ExpandLifecycleRules(l []interface{}) ([]*s3.LifecycleRule, error) {
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}

	var results []*s3.LifecycleRule

	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		result := &s3.LifecycleRule{}

		if v, ok := tfMap["abort_incomplete_multipart_upload"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			result.AbortIncompleteMultipartUpload = ExpandLifecycleRuleAbortIncompleteMultipartUpload(v[0].(map[string]interface{}))
		}

		if v, ok := tfMap["expiration"].([]interface{}); ok && len(v) > 0 {
			expiration, err := ExpandLifecycleRuleExpiration(v)
			if err != nil {
				return nil, err
			}
			result.Expiration = expiration
		}

		if v, ok := tfMap["filter"].([]interface{}); ok && len(v) > 0 {
			result.Filter = ExpandLifecycleRuleFilter(v)
		}

		if v, ok := tfMap["id"].(string); ok {
			result.ID = aws.String(v)
		}

		if v, ok := tfMap["noncurrent_version_expiration"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			result.NoncurrentVersionExpiration = ExpandLifecycleRuleNoncurrentVersionExpiration(v[0].(map[string]interface{}))
		}

		if v, ok := tfMap["noncurrent_version_transition"].(*schema.Set); ok && v.Len() > 0 {
			result.NoncurrentVersionTransitions = ExpandLifecycleRuleNoncurrentVersionTransitions(v.List())
		}

		if v, ok := tfMap["prefix"].(string); ok && result.Filter == nil {
			result.Prefix = aws.String(v)
		}

		if v, ok := tfMap["status"].(string); ok && v != "" {
			result.Status = aws.String(v)
		}

		if v, ok := tfMap["transition"].(*schema.Set); ok && v.Len() > 0 {
			transitions, err := ExpandLifecycleRuleTransitions(v.List())
			if err != nil {
				return nil, err
			}
			result.Transitions = transitions
		}

		results = append(results, result)
	}

	return results, nil
}

func ExpandMetrics(l []interface{}) *s3.Metrics {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.Metrics{}

	if v, ok := tfMap["event_threshold"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.EventThreshold = ExpandReplicationTimeValue(v)
	}

	if v, ok := tfMap["status"].(string); ok && v != "" {
		result.Status = aws.String(v)
	}

	return result
}

func ExpandReplicationRuleAndOperator(l []interface{}) *s3.ReplicationRuleAndOperator {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.ReplicationRuleAndOperator{}

	if v, ok := tfMap["prefix"].(string); ok && v != "" {
		result.Prefix = aws.String(v)
	}

	if v, ok := tfMap["tags"].(map[string]interface{}); ok && len(v) > 0 {
		tags := Tags(tftags.New(v).IgnoreAWS())
		if len(tags) > 0 {
			result.Tags = tags
		}
	}

	return result
}

func ExpandReplicationTime(l []interface{}) *s3.ReplicationTime {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.ReplicationTime{}

	if v, ok := tfMap["status"].(string); ok && v != "" {
		result.Status = aws.String(v)
	}

	if v, ok := tfMap["time"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.Time = ExpandReplicationTimeValue(v)
	}

	return result
}

func ExpandReplicationTimeValue(l []interface{}) *s3.ReplicationTimeValue {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.ReplicationTimeValue{}

	if v, ok := tfMap["minutes"].(int); ok {
		result.Minutes = aws.Int64(int64(v))
	}

	return result
}

func ExpandReplicaModifications(l []interface{}) *s3.ReplicaModifications {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.ReplicaModifications{}

	if v, ok := tfMap["status"].(string); ok && v != "" {
		result.Status = aws.String(v)
	}

	return result
}

func ExpandRules(l []interface{}) []*s3.ReplicationRule {
	var rules []*s3.ReplicationRule

	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}
		rule := &s3.ReplicationRule{}

		if v, ok := tfMap["delete_marker_replication"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			rule.DeleteMarkerReplication = ExpandDeleteMarkerReplication(v)
		}

		if v, ok := tfMap["destination"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			rule.Destination = ExpandDestination(v)
		}

		if v, ok := tfMap["existing_object_replication"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			rule.ExistingObjectReplication = ExpandExistingObjectReplication(v)
		}

		if v, ok := tfMap["id"].(string); ok && v != "" {
			rule.ID = aws.String(v)
		}

		if v, ok := tfMap["source_selection_criteria"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			rule.SourceSelectionCriteria = ExpandSourceSelectionCriteria(v)
		}

		if v, ok := tfMap["status"].(string); ok && v != "" {
			rule.Status = aws.String(v)
		}

		if v, ok := tfMap["filter"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
			// XML schema V2
			rule.Filter = ExpandFilter(v)
			rule.Priority = aws.Int64(int64(tfMap["priority"].(int)))
		} else {
			// XML schema V1
			rule.Prefix = aws.String(tfMap["prefix"].(string))
		}

		rules = append(rules, rule)
	}

	return rules
}

func ExpandSourceSelectionCriteria(l []interface{}) *s3.SourceSelectionCriteria {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.SourceSelectionCriteria{}

	if v, ok := tfMap["replica_modifications"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.ReplicaModifications = ExpandReplicaModifications(v)
	}

	if v, ok := tfMap["sse_kms_encrypted_objects"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		result.SseKmsEncryptedObjects = ExpandSseKmsEncryptedObjects(v)
	}

	return result
}

func ExpandSseKmsEncryptedObjects(l []interface{}) *s3.SseKmsEncryptedObjects {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.SseKmsEncryptedObjects{}

	if v, ok := tfMap["status"].(string); ok && v != "" {
		result.Status = aws.String(v)
	}

	return result
}

func ExpandTag(l []interface{}) *s3.Tag {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	result := &s3.Tag{}

	if v, ok := tfMap["key"].(string); ok && v != "" {
		result.Key = aws.String(v)
	}

	if v, ok := tfMap["value"].(string); ok && v != "" {
		result.Value = aws.String(v)
	}

	return result
}

func FlattenAccessControlTranslation(act *s3.AccessControlTranslation) []interface{} {
	if act == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if act.Owner != nil {
		m["owner"] = aws.StringValue(act.Owner)
	}

	return []interface{}{m}
}

func FlattenEncryptionConfiguration(ec *s3.EncryptionConfiguration) []interface{} {
	if ec == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if ec.ReplicaKmsKeyID != nil {
		m["replica_kms_key_id"] = aws.StringValue(ec.ReplicaKmsKeyID)
	}

	return []interface{}{m}
}

func FlattenDeleteMarkerReplication(dmr *s3.DeleteMarkerReplication) []interface{} {
	if dmr == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if dmr.Status != nil {
		m["status"] = aws.StringValue(dmr.Status)
	}

	return []interface{}{m}
}

func FlattenDestination(dest *s3.Destination) []interface{} {
	if dest == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if dest.AccessControlTranslation != nil {
		m["access_control_translation"] = FlattenAccessControlTranslation(dest.AccessControlTranslation)
	}

	if dest.Account != nil {
		m["account"] = aws.StringValue(dest.Account)
	}

	if dest.Bucket != nil {
		m["bucket"] = aws.StringValue(dest.Bucket)
	}

	if dest.EncryptionConfiguration != nil {
		m["encryption_configuration"] = FlattenEncryptionConfiguration(dest.EncryptionConfiguration)
	}

	if dest.Metrics != nil {
		m["metrics"] = FlattenMetrics(dest.Metrics)
	}

	if dest.ReplicationTime != nil {
		m["replication_time"] = FlattenReplicationTime(dest.ReplicationTime)
	}

	if dest.StorageClass != nil {
		m["storage_class"] = aws.StringValue(dest.StorageClass)
	}

	return []interface{}{m}
}

func FlattenExistingObjectReplication(eor *s3.ExistingObjectReplication) []interface{} {
	if eor == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if eor.Status != nil {
		m["status"] = aws.StringValue(eor.Status)
	}

	return []interface{}{m}
}

func FlattenFilter(filter *s3.ReplicationRuleFilter) []interface{} {
	if filter == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if filter.And != nil {
		m["and"] = FlattenReplicationRuleAndOperator(filter.And)
	}

	if filter.Prefix != nil {
		m["prefix"] = aws.StringValue(filter.Prefix)
	}

	if filter.Tag != nil {
		tag := KeyValueTags([]*s3.Tag{filter.Tag}).IgnoreAWS().Map()
		m["tag"] = []interface{}{tag}
	}

	return []interface{}{m}
}

func FlattenLifecycleRules(rules []*s3.LifecycleRule) []interface{} {
	if len(rules) == 0 {
		return []interface{}{}
	}

	var results []interface{}

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		m := make(map[string]interface{})

		if rule.AbortIncompleteMultipartUpload != nil {
			m["abort_incomplete_multipart_upload"] = FlattenLifecycleRuleAbortIncompleteMultipartUpload(rule.AbortIncompleteMultipartUpload)
		}

		if rule.Expiration != nil {
			m["expiration"] = FlattenLifecycleRuleExpiration(rule.Expiration)
		}

		if rule.Filter != nil {
			m["filter"] = FlattenLifecycleRuleFilter(rule.Filter)
		}

		if rule.ID != nil {
			m["id"] = aws.StringValue(rule.ID)
		}

		if rule.NoncurrentVersionExpiration != nil {
			m["noncurrent_version_expiration"] = FlattenLifecycleRuleNoncurrentVersionExpiration(rule.NoncurrentVersionExpiration)
		}

		if rule.NoncurrentVersionTransitions != nil {
			m["noncurrent_version_transition"] = FlattenLifecycleRuleNoncurrentVersionTransitions(rule.NoncurrentVersionTransitions)
		}

		if rule.Prefix != nil {
			m["prefix"] = aws.StringValue(rule.Prefix)
		}

		if rule.Status != nil {
			m["status"] = aws.StringValue(rule.Status)
		}

		if rule.Transitions != nil {
			m["transition"] = FlattenLifecycleRuleTransitions(rule.Transitions)
		}

		results = append(results, m)
	}

	return results
}

func FlattenLifecycleRuleAbortIncompleteMultipartUpload(u *s3.AbortIncompleteMultipartUpload) []interface{} {
	if u == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if u.DaysAfterInitiation != nil {
		m["days_after_initiation"] = int(aws.Int64Value(u.DaysAfterInitiation))
	}

	return []interface{}{m}
}

func FlattenLifecycleRuleExpiration(expiration *s3.LifecycleExpiration) []interface{} {
	if expiration == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if expiration.Days != nil {
		m["days"] = int(aws.Int64Value(expiration.Days))
	}

	if expiration.Date != nil {
		m["date"] = expiration.Date.Format(time.RFC3339)
	}

	if expiration.ExpiredObjectDeleteMarker != nil {
		m["expired_object_delete_marker"] = aws.BoolValue(expiration.ExpiredObjectDeleteMarker)
	}

	return []interface{}{m}
}

func FlattenLifecycleRuleFilter(filter *s3.LifecycleRuleFilter) []interface{} {
	if filter == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if filter.And != nil {
		m["and"] = FlattenLifecycleRuleFilterAndOperator(filter.And)
	}

	if filter.ObjectSizeGreaterThan != nil {
		m["object_size_greater_than"] = int(aws.Int64Value(filter.ObjectSizeGreaterThan))
	}

	if filter.ObjectSizeLessThan != nil {
		m["object_size_less_than"] = int(aws.Int64Value(filter.ObjectSizeLessThan))
	}

	if filter.Prefix != nil {
		m["prefix"] = aws.StringValue(filter.Prefix)
	}

	if filter.Tag != nil {
		m["tag"] = FlattenLifecycleRuleFilterTag(filter.Tag)
	}

	return []interface{}{m}
}

func FlattenLifecycleRuleFilterAndOperator(andOp *s3.LifecycleRuleAndOperator) []interface{} {
	if andOp == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if andOp.ObjectSizeGreaterThan != nil {
		m["object_size_greater_than"] = int(aws.Int64Value(andOp.ObjectSizeGreaterThan))
	}

	if andOp.ObjectSizeLessThan != nil {
		m["object_size_less_than"] = int(aws.Int64Value(andOp.ObjectSizeLessThan))
	}

	if andOp.Prefix != nil {
		m["prefix"] = aws.StringValue(andOp.Prefix)
	}

	if andOp.Tags != nil {
		m["tags"] = KeyValueTags(andOp.Tags).IgnoreAWS().Map()
	}

	return []interface{}{m}
}

func FlattenLifecycleRuleFilterTag(tag *s3.Tag) []interface{} {
	if tag == nil {
		return []interface{}{}
	}

	t := KeyValueTags([]*s3.Tag{tag}).IgnoreAWS().Map()

	return []interface{}{t}
}

func FlattenLifecycleRuleNoncurrentVersionExpiration(expiration *s3.NoncurrentVersionExpiration) []interface{} {
	if expiration == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if expiration.NewerNoncurrentVersions != nil {
		m["newer_noncurrent_versions"] = int(aws.Int64Value(expiration.NewerNoncurrentVersions))
	}

	if expiration.NoncurrentDays != nil {
		m["noncurrent_days"] = int(aws.Int64Value(expiration.NoncurrentDays))
	}

	return []interface{}{m}
}

func FlattenLifecycleRuleNoncurrentVersionTransitions(transitions []*s3.NoncurrentVersionTransition) []interface{} {
	if len(transitions) == 0 {
		return []interface{}{}
	}

	var results []interface{}

	for _, transition := range transitions {
		if transition == nil {
			continue
		}

		m := make(map[string]interface{})

		if transition.NewerNoncurrentVersions != nil {
			m["newer_noncurrent_versions"] = int(aws.Int64Value(transition.NewerNoncurrentVersions))
		}

		if transition.NoncurrentDays != nil {
			m["noncurrent_days"] = int(aws.Int64Value(transition.NoncurrentDays))
		}

		if transition.StorageClass != nil {
			m["storage_class"] = aws.StringValue(transition.StorageClass)
		}

		results = append(results, m)
	}

	return results
}

func FlattenLifecycleRuleTransitions(transitions []*s3.Transition) []interface{} {
	if len(transitions) == 0 {
		return []interface{}{}
	}

	var results []interface{}

	for _, transition := range transitions {
		if transition == nil {
			continue
		}

		m := make(map[string]interface{})

		if transition.Date != nil {
			m["date"] = transition.Date.Format(time.RFC3339)
		}

		if transition.Days != nil {
			m["days"] = int(aws.Int64Value(transition.Days))
		}

		if transition.StorageClass != nil {
			m["storage_class"] = aws.StringValue(transition.StorageClass)
		}

		results = append(results, m)
	}

	return results
}

func FlattenMetrics(metrics *s3.Metrics) []interface{} {
	if metrics == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if metrics.EventThreshold != nil {
		m["event_threshold"] = FlattenReplicationTimeValue(metrics.EventThreshold)
	}

	if metrics.Status != nil {
		m["status"] = aws.StringValue(metrics.Status)
	}

	return []interface{}{m}
}

func FlattenReplicationTime(rt *s3.ReplicationTime) []interface{} {
	if rt == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if rt.Status != nil {
		m["status"] = aws.StringValue(rt.Status)
	}

	if rt.Time != nil {
		m["time"] = FlattenReplicationTimeValue(rt.Time)
	}

	return []interface{}{m}

}

func FlattenReplicationTimeValue(rtv *s3.ReplicationTimeValue) []interface{} {
	if rtv == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if rtv.Minutes != nil {
		m["minutes"] = int(aws.Int64Value(rtv.Minutes))
	}

	return []interface{}{m}
}

func FlattenRules(rules []*s3.ReplicationRule) []interface{} {
	if len(rules) == 0 {
		return []interface{}{}
	}

	var results []interface{}

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		m := make(map[string]interface{})

		if rule.DeleteMarkerReplication != nil {
			m["delete_marker_replication"] = FlattenDeleteMarkerReplication(rule.DeleteMarkerReplication)
		}

		if rule.Destination != nil {
			m["destination"] = FlattenDestination(rule.Destination)
		}

		if rule.ExistingObjectReplication != nil {
			m["existing_object_replication"] = FlattenExistingObjectReplication(rule.ExistingObjectReplication)
		}

		if rule.Filter != nil {
			m["filter"] = FlattenFilter(rule.Filter)
		}

		if rule.ID != nil {
			m["id"] = aws.StringValue(rule.ID)
		}

		if rule.Prefix != nil {
			m["prefix"] = aws.StringValue(rule.Prefix)
		}

		if rule.Priority != nil {
			m["priority"] = int(aws.Int64Value(rule.Priority))
		}

		if rule.SourceSelectionCriteria != nil {
			m["source_selection_criteria"] = FlattenSourceSelectionCriteria(rule.SourceSelectionCriteria)
		}

		if rule.Status != nil {
			m["status"] = aws.StringValue(rule.Status)
		}

		results = append(results, m)
	}

	return results
}

func FlattenReplicaModifications(rc *s3.ReplicaModifications) []interface{} {
	if rc == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if rc.Status != nil {
		m["status"] = aws.StringValue(rc.Status)
	}

	return []interface{}{m}
}

func FlattenReplicationRuleAndOperator(op *s3.ReplicationRuleAndOperator) []interface{} {
	if op == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if op.Prefix != nil {
		m["prefix"] = aws.StringValue(op.Prefix)
	}

	if op.Tags != nil {
		m["tags"] = KeyValueTags(op.Tags).IgnoreAWS().Map()
	}

	return []interface{}{m}

}

func FlattenSourceSelectionCriteria(ssc *s3.SourceSelectionCriteria) []interface{} {
	if ssc == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if ssc.ReplicaModifications != nil {
		m["replica_modifications"] = FlattenReplicaModifications(ssc.ReplicaModifications)
	}

	if ssc.SseKmsEncryptedObjects != nil {
		m["sse_kms_encrypted_objects"] = FlattenSseKmsEncryptedObjects(ssc.SseKmsEncryptedObjects)
	}

	return []interface{}{m}
}

func FlattenSseKmsEncryptedObjects(objects *s3.SseKmsEncryptedObjects) []interface{} {
	if objects == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if objects.Status != nil {
		m["status"] = aws.StringValue(objects.Status)
	}

	return []interface{}{m}
}
