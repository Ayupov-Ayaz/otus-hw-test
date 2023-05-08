CREATE TABLE notifications (
  event_id int(11) NOT NULL,
  before_start_notice_sec INT NOT NULL,
  PRIMARY KEY (event_id, before_start_notice_sec),
  FOREIGN KEY (event_id) REFERENCES events (id)
);