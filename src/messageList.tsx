import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import {
  makeStyles,
  useTheme,
  Card,
  CardContent,
  Typography,
} from '@material-ui/core';
// eslint-disable-next-line no-unused-vars
import Box from '@material-ui/core/Box';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import SwipeableViews from 'react-swipeable-views';
import { Message, User } from './datatypes';
import EditMessagePostForm from './editForm';
import DeleteMessageDialog from './deleteDialog';
import DownloadMessageDialog from './downloadDialog';
import AudioMessagePostForm from './audioMessageForm';
import MessagePostForm from './messageForm';
import DataAnalysisPage from './dataAnalysisPage';

const useStylesCard = makeStyles({
  root: {
    minWidth: 200,
    maxWidth: 1000,
  },
  bullet: {
    display: 'inline-block',
    margin: '0 2px',
    transform: 'scale(0.8)',
  },
  title: {
    fontSize: 14,
  },
});

const useTabStyles = makeStyles((theme) => ({
  tab: {
    backgroundColor: theme.palette.background.paper,
    width: 500,
  },
  analysis: {
    width: 500,
  },
}));

interface TabPanelProps {
  children: React.ReactNode;
  index: any;
  value: any;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`full-width-tabpanel-${index}`}
      aria-labelledby={`full-width-tab-${index}`}
    >
      {value === index && (
        <Box p={3}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

function GetMessage({ forceUpdate }) {
  const classes = useStylesCard();
  const [data, setData] = useState<Message[]>([]);
  // ユーザー情報を取得
  const [userData, setUserData] = useState<User>({ id: 0, name: '' });

  useEffect(() => {
    // ルート /message に対して GETリクエストを送る
    // 帰ってきたものをjsonにしてuseStateに突っ込む
    fetch(`${window.location.pathname}/message`)
      .then((res) => res.json())
      .then(setData);
  }, [forceUpdate]);

  useEffect(() => {
    fetch('/user')
      .then((res) => res.json())
      .then(setUserData);
  }, []);

  return (
    // タグが複数できる場合は何らかのタグで全体を囲う
    <div>
      <DownloadMessageDialog
        targetMessage={data
          .reduce(
            (prev, current) =>
              // 書き込み順で文字列にしている
              `[${current.addedBy.name}]\n${current.message}\n\n${prev}`,
            ''
          )
          .toString()}
        title={document.title}
      />

      {data.map((item) => (
        <Card className={classes.root} key={item.id}>
          <CardContent>
            <Typography
              className={classes.title}
              color="textSecondary"
              gutterBottom
              align="left"
            >
              {item.addedBy.name}
              <EditMessagePostForm
                prevMessage={item.message}
                id={item.id.toString()}
                isHidden={userData.id !== item.addedBy.id}
              />
              <DeleteMessageDialog
                targetMessage={item.message}
                id={item.id.toString()}
                isHidden={userData.id !== item.addedBy.id}
              />
            </Typography>
            <Typography variant="body2" component="p" align="left">
              {item.message}
            </Typography>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

GetMessage.propTypes = {
  // このランダム値を変更することで、強制的にサーバーからメッセージを取得させ、最新の情報を入手させる
  forceUpdate: PropTypes.number,
};

GetMessage.defaultProps = {
  forceUpdate: Math.random(),
};

export default function MessageList() {
  const classes = useTabStyles();
  const theme = useTheme();
  const [value, setValue] = React.useState(0);
  const [randomValue, setRandomValue] = useState<number>(Math.random());

  const onMessageAdded = () => {
    // フォームによってメッセージが追加されたら、メッセージ一覧を更新する
    setRandomValue(Math.random());
  };

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setValue(newValue);
  };

  const handleChangeIndex = (index: number) => {
    setValue(index);
  };

  return (
    <>
      <DataAnalysisPage />
      <div className={classes.tab}>
        <AppBar position="static" color="default">
          <Tabs
            value={value}
            onChange={handleChange}
            indicatorColor="primary"
            textColor="primary"
            variant="fullWidth"
            aria-label="full width tabs example"
          >
            <Tab
              label="音声入力"
              id="full-width-tab-0"
              aria-controls="full-width-tabpanel-0"
            />
            <Tab
              label="キーボード入力"
              id="full-width-tab-1"
              aria-controls="full-width-tabpanel-1"
            />
          </Tabs>
        </AppBar>
        <SwipeableViews
          axis={theme.direction === 'rtl' ? 'x-reverse' : 'x'}
          index={value}
          onChangeIndex={handleChangeIndex}
        >
          <TabPanel value={value} index={0}>
            <AudioMessagePostForm onSubmitSuccessful={onMessageAdded} />
          </TabPanel>
          <TabPanel value={value} index={1}>
            <MessagePostForm onSubmitSuccessful={onMessageAdded} />
          </TabPanel>
        </SwipeableViews>
      </div>
      <GetMessage forceUpdate={randomValue} />
    </>
  );
}
