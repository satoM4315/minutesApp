import React, { useState } from 'react';
import PropTypes from 'prop-types';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';

// メッセージ更新のAPIへのURL
// eslint-disable-next-line @typescript-eslint/naming-convention
const API_URL_UPDATE_MESSAGE = '/update_message';

export default function EditMessagePostForm(props) {
  const { onSubmitSuccessful, prevMessage, id, isHidden } = props;

  const [message, setMessage] = React.useState<string>(prevMessage);

  const [open, setOpen] = React.useState<boolean>(false);
  // サーバがへメッセージ追加のリクエストを処理中ならtrue、でないならfalseの状態
  const [working, setWorking] = useState<boolean>(false);

  const handleUpdate = async (event: React.FormEvent) => {
    // FIXME もしかしたら、非同期なため、これが効く前にボタンをクリックできるかもしれない
    setWorking(true);
    try {
      // ページが更新されないようにする
      event.preventDefault();

      // Reactのハンドラはasyncにできる
      const res = await fetch(API_URL_UPDATE_MESSAGE, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        // 相応しくないかも
        // same-originを使うべき？
        credentials: 'include',
        body: JSON.stringify({ id, message }),
      });
      const obj = await res.json();
      if ('error' in obj) {
        // サーバーからエラーが返却された
        throw new Error(
          `An error occurred on querying ${API_URL_UPDATE_MESSAGE}, the response included error message: ${obj.error}`
        );
      }
      if (!('success' in obj)) {
        // サーバーからsuccessメンバが含まれたJSONが帰るはずだが、見当たらなかった
        throw new Error(
          `An response from ${API_URL_UPDATE_MESSAGE} unexpectedly did not have 'success' member`
        );
      }
      if (obj.success !== true) {
        throw new Error(
          `An response from ${API_URL_UPDATE_MESSAGE} returned non true value as 'success' member`
        );
      }
      // 要求は成功
      // リスナ関数を呼ぶ
      onSubmitSuccessful();
    } finally {
      setWorking(false);
      setMessage('');
      setOpen(false);
    }
  };

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  if (isHidden) {
    return null;
  }
  return (
    <span>
      <IconButton
        edge="end"
        aria-label="edit"
        onClick={handleClickOpen}
        data-num="100"
        disabled={working}
      >
        <EditIcon fontSize="small" />
      </IconButton>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="form-dialog-title"
        fullWidth
      >
        <DialogTitle id="form-dialog-title">編集</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            type="text"
            fullWidth
            multiline
            value={message}
            onChange={(event) => setMessage(event.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="primary">
            キャンセル
          </Button>
          <Button onClick={handleUpdate} color="primary">
            完了
          </Button>
        </DialogActions>
      </Dialog>
    </span>
  );
}

EditMessagePostForm.propTypes = {
  // メッセージの更新が正常に完了したら呼ばれる関数
  onSubmitSuccessful: PropTypes.func,
  prevMessage: PropTypes.string,
  id: PropTypes.string,
  isHidden: PropTypes.bool,
};

EditMessagePostForm.defaultProps = {
  onSubmitSuccessful: () => {
    window.location.href = window.location.pathname;
  },
  prevMessage: '',
  id: '',
  isHidden: true,
};
