package kafkaclient

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Cluster struct {
	cluster sarama.ClusterAdmin
}

// NewCluster creates a new ClusterAdmin client from given brokers and sarama config.
// If cfg is nil, a default config will be used.
func NewCluster(brokers []string, cfg *sarama.Config) (*Cluster, error) {
	if cfg == nil {
		cfg = sarama.NewConfig()
		cfg.Version = sarama.V2_8_0_0 // or your desired default version
	}

	cluster, err := sarama.NewClusterAdmin(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("NewClusterAdmin error: %w", err)
	}

	return &Cluster{
		cluster: cluster,
	}, nil
}

// Topics returns a list of all topic names in the cluster.
func (c *Cluster) Topics() ([]string, error) {
	topicsMap, err := c.cluster.ListTopics()
	if err != nil {
		return nil, fmt.Errorf("cluster.ListTopics error: %w", err)
	}

	var topics []string
	for topic := range topicsMap {
		topics = append(topics, topic)
	}

	return topics, nil
}

// GetTopic fetches the details of a given topic.
// Returns the TopicDetail, a boolean indicating if the topic exists, and an error if any.
func (c *Cluster) GetTopic(topic string) (sarama.TopicDetail, bool, error) {
	topicsMap, err := c.cluster.ListTopics()
	if err != nil {
		return sarama.TopicDetail{}, false, fmt.Errorf("cluster.ListTopics error: %w", err)
	}

	detail, exists := topicsMap[topic]
	if !exists {
		return sarama.TopicDetail{}, false, nil
	}

	return detail, true, nil
}

// CreateTopic creates a topic with the given detail.
// If validateOnly is true, it only validates the topic without creating it.
func (c *Cluster) CreateTopic(topic string, detail *sarama.TopicDetail, validateOnly bool) error {
	exists, err := func() (bool, error) {
		_, ok, err := c.GetTopic(topic)
		return ok, err
	}()
	if err != nil {
		return err
	}

	if exists {
		return nil // topic already exists, no error
	}

	if err := c.cluster.CreateTopic(topic, detail, validateOnly); err != nil {
		return fmt.Errorf("cluster.CreateTopic error: %w", err)
	}

	return nil
}

// DeleteTopic deletes the specified topic.
func (c *Cluster) DeleteTopic(topic string) error {
	if err := c.cluster.DeleteTopic(topic); err != nil {
		return fmt.Errorf("cluster.DeleteTopic error: %w", err)
	}
	return nil
}

// ConsumerGroups returns all consumer group IDs of type "consumer".
func (c *Cluster) ConsumerGroups() ([]string, error) {
	groups, err := c.cluster.ListConsumerGroups()
	if err != nil {
		return nil, fmt.Errorf("cluster.ListConsumerGroups error: %w", err)
	}

	var consumers []string
	for groupID, protocolType := range groups {
		// Usually, protocolType "consumer" means a standard consumer group
		if protocolType == "consumer" {
			consumers = append(consumers, groupID)
		}
	}

	return consumers, nil
}
