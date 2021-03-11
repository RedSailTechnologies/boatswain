package run

import (
	"errors"
	"regexp"
)

type statuses struct {
	history map[string]Status
	last    Status
	overall Status
}

func (s *statuses) addStatus(step string, status Status) {
	if s.history == nil {
		s.history = make(map[string]Status)
	}

	s.history[step] = status
	s.last = status
	if status != Skipped {
		s.overall = status
	}
}

func (s *statuses) shouldExecute(c string) (bool, error) {
	cond, err := parseCondition(c)
	if err != nil {
		return false, err
	}

	var status Status
	var ok bool
	if cond.cond == "always" {
		return true, nil
	} else if cond.step == "" {
		status = s.last
	} else if status, ok = s.history[cond.step]; !ok {
		return false, errors.New("step for condition not found")
	}

	switch cond.cond {
	case "any":
		return true, nil
	case "failed":
		if status == Failed {
			return true, nil
		}
	case "succeeded":
		if status == Succeeded {
			return true, nil
		}
	case "skipped":
		if status == Skipped {
			return true, nil
		}
	}
	return false, nil
}

type condition struct {
	cond string
	step string
}

func parseCondition(c string) (*condition, error) {
	re := regexp.MustCompile(`^(always|any|failed|skipped|succeeded)\((.*)\)$`)
	m := re.FindStringSubmatch(c)
	if len(m) != 3 || m[0] != c {
		return nil, errors.New("could not parse condition")
	}
	return &condition{cond: m[1], step: m[2]}, nil
}
