package main

type Writer interface {
	GetCanonicalName() string
	Write(*TradeRow) error
	Close() error
}

type WriterChannel struct {
	writer Writer
	trades chan TradeRow
}

// Stores all writers
type Writers struct {
	_Writers        []Writer
	_WriterChannels []WriterChannel
	CloseChannel    chan bool
}

// Register a writer
func (writers *Writers) RegisterWriter(writer Writer) {
	writers._Writers = append(writers._Writers, writer)
}

// Writes to all writers in a separate go routine
func (writers *Writers) Write(tradeRows []TradeRow) error {
	go func() {
		select {
		case <-writers.CloseChannel:
			return
		default:
			{
				for _, writerChannel := range writers._WriterChannels {
					for _, tradeRow := range tradeRows {
						if writerChannel.trades != nil {
							writerChannel.trades <- tradeRow
						} //else - lost logger
					}
				}
			}
		}
	}()
	return nil
}

// Start a go routine for each writer and create a struct of channels to write on
func (writers *Writers) Listen(closeChannel chan bool) {
	var _writerChannels []WriterChannel
	// make the channel and create the WriterChannel
	for _, _writer := range writers._Writers {
		var _writerChannel = WriterChannel{
			writer: _writer,
			trades: make(chan TradeRow),
		}
		_writerChannels = append(_writerChannels, _writerChannel)
		go func(writer Writer, writerChannel WriterChannel) {
			defer func() {
				if r := recover(); r != nil {
					LOGGER.Error.Println("Error occured while writing [", writer.GetCanonicalName(), "]: ", r)
				}
			}()
			for {
				select {
				case tradeRow := <-writerChannel.trades:
					writer.Write(&tradeRow)
				case <-closeChannel:
					close(writerChannel.trades)
					writerChannel.trades = nil
					writer.Close()
					return
				}
			}
		}(_writer, _writerChannel)
	}
	writers._WriterChannels = _writerChannels
	writers.CloseChannel = closeChannel
}

func (writers *Writers) GetWriters() []Writer {
	return writers._Writers
}

var WRITERS = Writers{}
