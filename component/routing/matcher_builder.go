/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2022-2025, daeuniverse Organization <dae@v2raya.org>
 */

package routing

import (
	"fmt"
	"github.com/daeuniverse/dae/common/consts"
	"github.com/daeuniverse/dae/pkg/config_parser"
	"github.com/sirupsen/logrus"
	"strconv"
)

type DomainSet struct {
	Key       consts.RoutingDomainKey
	RuleIndex int
	Domains   []string
}

type Outbound struct {
	Name string
	Mark uint32
	Must bool
}

type RulesBuilder struct {
	log     *logrus.Logger
	parsers map[string]FunctionParser
}

func NewRulesBuilder(log *logrus.Logger) *RulesBuilder {
	return &RulesBuilder{
		log:     log,
		parsers: make(map[string]FunctionParser),
	}
}

func (b *RulesBuilder) RegisterFunctionParser(funcName string, parser FunctionParser) {
	b.parsers[funcName] = parser
}

func (b *RulesBuilder) Apply(rules []*config_parser.RoutingRule) (err error) {
	for _, rule := range rules {
		b.log.Debugln("[rule]", rule.String(true, false, false))
		outbound, err := ParseOutbound(&rule.Outbound)
		if err != nil {
			return err
		}

		// rule is like: domain(domain:baidu.com) && port(443) -> proxy
		for iFunc, f := range rule.AndFunctions {
			// f is like: domain(domain:baidu.com)
			functionParser, ok := b.parsers[f.Name]
			if !ok {
				return fmt.Errorf("unknown function: %v", f.Name)
			}
			paramValueGroups := groupParamValuesByKey(f.Params)
			for jMatchSet, group := range paramValueGroups {
				key := group.key
				// Preprocess the outbound.
				overrideOutbound := &Outbound{
					Name: consts.OutboundLogicalOr.String(),
					Mark: outbound.Mark,
					Must: outbound.Must,
				}
				if jMatchSet == len(paramValueGroups)-1 {
					overrideOutbound.Name = consts.OutboundLogicalAnd.String()
					if iFunc == len(rule.AndFunctions)-1 {
						overrideOutbound.Name = outbound.Name
					}
				}

				{
					// Debug
					symNot := ""
					if f.Not {
						symNot = "!"
					}
					b.log.Debugf("\t%v%v(%v) -> %v", symNot, f.Name, key, overrideOutbound.Name)
				}

				if err = functionParser(b.log, f, key, group.values, overrideOutbound); err != nil {
					return fmt.Errorf("failed to parse '%v': %w", f.String(false, false, false), err)
				}
			}
		}
	}
	return nil
}

type paramValueGroup struct {
	key    string
	values []string
}

func groupParamValuesByKey(params []*config_parser.Param) []paramValueGroup {
	var groups []paramValueGroup
	sorted := true
	for i := 0; i < len(params); {
		j := i + 1
		for j < len(params) && params[j].Key == params[i].Key {
			j++
		}
		values := make([]string, j-i)
		for k := i; k < j; k++ {
			values[k-i] = params[k].Val
		}
		if len(groups) > 0 && params[i].Key < groups[len(groups)-1].key {
			sorted = false
		}
		groups = append(groups, paramValueGroup{key: params[i].Key, values: values})
		i = j
	}

	if sorted {
		return groups
	}
	keyToGroup := make(map[string]int, len(groups))
	compacted := groups[:0]
	for _, group := range groups {
		if i, ok := keyToGroup[group.key]; ok {
			compacted[i].values = append(compacted[i].values, group.values...)
			continue
		}
		keyToGroup[group.key] = len(compacted)
		compacted = append(compacted, group)
	}
	return compacted
}

func ParseOutbound(rawOutbound *config_parser.Function) (outbound *Outbound, err error) {
	outbound = &Outbound{
		Name: rawOutbound.Name,
		Mark: 0,
		Must: false,
	}
	for _, p := range rawOutbound.Params {
		switch p.Key {
		case consts.OutboundParam_Mark:
			var _mark uint64
			_mark, err = strconv.ParseUint(p.Val, 0, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mark: %v", err)
			}
			outbound.Mark = uint32(_mark)
		case "":
			if p.Val == "must" {
				outbound.Must = true
			} else {
				return nil, fmt.Errorf("unknown outbound param: %v", p.Val)
			}
		default:
			return nil, fmt.Errorf("unknown outbound param key: %v", p.Key)
		}
	}
	return outbound, nil
}
