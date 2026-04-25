"use client";

import { api } from "@/lib/api";
import type { EventItem, JobItem, MetricsSummary } from "@/lib/types";
import {
  Button,
  Card,
  Col,
  Input,
  Modal,
  Row,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from "antd";
import type { ColumnsType } from "antd/es/table";
import Link from "next/link";
import { useEffect, useState } from "react";

function statusColor(status: string) {
  switch (status) {
    case "succeeded":
    case "planned":
      return "green";
    case "pending":
      return "blue";
    case "processing":
      return "orange";
    case "dead_lettered":
    case "failed":
      return "red";
    default:
      return "default";
  }
}

export default function OverviewPage() {
  const [metrics, setMetrics] = useState<MetricsSummary | null>(null);
  const [events, setEvents] = useState<EventItem[]>([]);
  const [jobs, setJobs] = useState<JobItem[]>([]);

  const [demoOpen, setDemoOpen] = useState(false);
  const [demoKey, setDemoKey] = useState("");
  const [demoLoading, setDemoLoading] = useState(false);

  async function loadDashboard() {
    const [metricsRes, eventsRes, jobsRes] = await Promise.all([
      api.getMetrics(),
      api.getEvents(),
      api.getJobs(),
    ]);

    setMetrics(metricsRes.data);
    setEvents(eventsRes.data.slice(0, 5));
    setJobs(jobsRes.data.slice(0, 5));
  }

  useEffect(() => {
    loadDashboard().catch(() => {
      message.error("Failed to load dashboard data");
    });
  }, []);

  async function generateDemoEvents() {
    if (!demoKey.trim()) {
      message.warning("Please enter demo key");
      return;
    }

    setDemoLoading(true);

    try {
      await api.createDemoEvents(10, demoKey.trim());
      message.success("Demo events created");
      setDemoOpen(false);
      setDemoKey("");
      await loadDashboard();
    } catch {
      message.error("Failed to create demo events");
    } finally {
      setDemoLoading(false);
    }
  }

  const eventColumns: ColumnsType<EventItem> = [
    {
      title: "Type",
      dataIndex: "event_type",
    },
    {
      title: "Source",
      dataIndex: "source",
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      render: (v) => new Date(v).toLocaleString(),
    },
  ];

  const jobColumns: ColumnsType<JobItem> = [
    {
      title: "Type",
      dataIndex: "job_type",
    },
    {
      title: "Channel",
      dataIndex: "channel",
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
    },
    {
      title: "Executions",
      render: (_, row) => `${row.attempts} total`,
    },
    {
      title: "Action",
      render: (_, row) => <Link href={`/jobs/${row.id}`}>View</Link>,
    },
  ];

  if (!metrics) {
    return <Typography.Text>Loading...</Typography.Text>;
  }

  return (
    <Space orientation="vertical" size="large" style={{ width: "100%" }}>
      <Space
        align="center"
        wrap
        style={{ width: "100%", justifyContent: "space-between" }}
      >
        <Typography.Title level={2} style={{ margin: 0 }}>
          Overview
        </Typography.Title>

        <Button type="primary" onClick={() => setDemoOpen(true)}>
          Generate Demo Events
        </Button>
      </Space>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card title="Total Events">{metrics.total_events}</Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card title="Pending Jobs">{metrics.pending_jobs}</Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card title="Succeeded Jobs">{metrics.succeeded_jobs}</Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card title="Dead Letters">{metrics.dead_lettered_jobs}</Card>
        </Col>

        <Col xs={24} sm={12} lg={6}>
          <Card title="Total Attempts">{metrics.total_attempts}</Card>
        </Col>
      </Row>

      <Card title="Recent Jobs">
        <Table
          rowKey="id"
          columns={jobColumns}
          dataSource={jobs}
          pagination={false}
          scroll={{ x: "max-content" }}
        />
      </Card>

      <Card title="Recent Events">
        <Table
          rowKey="id"
          columns={eventColumns}
          dataSource={events}
          pagination={false}
          scroll={{ x: "max-content" }}
        />
      </Card>

      <Modal
        title="Generate Demo Events"
        open={demoOpen}
        onCancel={() => setDemoOpen(false)}
        onOk={generateDemoEvents}
        confirmLoading={demoLoading}
        okText="Generate"
      >
        <Space orientation="vertical" style={{ width: "100%" }}>
          <Typography.Text>
            Enter the demo key to generate sample events for the worker
            pipeline.
          </Typography.Text>

          <Input.Password
            placeholder="Demo key"
            value={demoKey}
            onChange={(e) => setDemoKey(e.target.value)}
            onPressEnter={generateDemoEvents}
          />
        </Space>
      </Modal>
    </Space>
  );
}
