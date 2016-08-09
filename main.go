package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

func main() {
	var (
		discovery  = flag.String("discovery", "", "type of discovery. Only ELB supported right now")
		awsRegion  = flag.String("aws-region", "eu-central-1", "AWS region")
		outputType = flag.String("output", "json", "output type, one of (json|query)")
		query      = flag.String("query", "", "template query")
	)
	flag.Parse()

	printer := func(val interface{}) error {
		switch *outputType {
		case "query":
			return template.Must(template.New("r").Parse(*query)).Execute(os.Stdout, val)
		default:
			out := struct {
				Data interface{} `json:"data"`
			}{val}
			return json.NewEncoder(os.Stdout).Encode(out)
		}
	}

	switch *discovery {
	case "ELB":
		err := getAllElasticLoadBalancers(*awsRegion, printer)
		if err != nil {
			log.Printf("Could not descibe load balancers: %v", err)
		}

	default:
		log.Printf("discovery type %s not supported", *discovery)
	}
}

func getAllElasticLoadBalancers(awsRegion string, printer func(val interface{}) error) error {
	svc := elb.New(session.New(), aws.NewConfig().WithRegion(awsRegion))
	params := &elb.DescribeLoadBalancersInput{}
	resp, err := svc.DescribeLoadBalancers(params)

	if err != nil {
		return fmt.Errorf("reading ELBs in region %q :%v", awsRegion, err)
	}

	elbs := make([]string, 0, len(resp.LoadBalancerDescriptions))

	for _, elb := range resp.LoadBalancerDescriptions {
		elbs = append(elbs, *elb.LoadBalancerName)
	}

	return printer(elbs)
}
