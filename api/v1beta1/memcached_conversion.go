package v1beta1

import (
	"fmt"
	"strings"

	v1 "github.com/example-inc/memcached-operator/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *Memcached) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.Memcached)

	sched := src.Spec.Schedule
	scheduleParts := []string{"*", "*", "*", "*", "*"}
	if sched.Minute != nil {
		scheduleParts[0] = string(*sched.Minute)
	}
	if sched.Hour != nil {
		scheduleParts[1] = string(*sched.Hour)
	}
	if sched.DayOfMonth != nil {
		scheduleParts[2] = string(*sched.DayOfMonth)
	}
	if sched.Month != nil {
		scheduleParts[3] = string(*sched.Month)
	}
	if sched.DayOfWeek != nil {
		scheduleParts[4] = string(*sched.DayOfWeek)
	}
	dst.Spec.Schedule = strings.Join(scheduleParts, " ")

	/*
		The rest of the conversion is pretty rote.
	*/
	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Size = src.Spec.Size
	// +kubebuilder:docs-gen:collapse=rote conversion
	return nil
}

func (dst *Memcached) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.Memcached)

	schedParts := strings.Split(src.Spec.Schedule, " ")
	if len(schedParts) != 5 {
		return fmt.Errorf("invalid schedule: not a standard 5-field schedule")
	}
	partIfNeeded := func(raw string) *CronField {
		if raw == "*" {
			return nil
		}
		part := CronField(raw)
		return &part
	}
	dst.Spec.Schedule.Minute = partIfNeeded(schedParts[0])
	dst.Spec.Schedule.Hour = partIfNeeded(schedParts[1])
	dst.Spec.Schedule.DayOfMonth = partIfNeeded(schedParts[2])
	dst.Spec.Schedule.Month = partIfNeeded(schedParts[3])
	dst.Spec.Schedule.DayOfWeek = partIfNeeded(schedParts[4])

	/*
		The rest of the conversion is pretty rote.
	*/
	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Size = src.Spec.Size
	// +kubebuilder:docs-gen:collapse=rote conversion
	return nil
}
